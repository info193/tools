package mq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 定义RabbitMQ对象
type RabbitMQ struct {
	Cfg       *RabbitMqConfig            // 配置
	Consumers map[*BusinessConfig]Handle // 消费者
	Lock      sync.Mutex                 // 锁
	retryCnfs map[int]int64
	retryNum  int
	done      chan struct{} // 优雅退出信号

	// 单例连接
	conn *amqp.Connection
	mu   sync.Mutex
}

// getConnection 获取单例连接，断线自动重建
func (r *RabbitMQ) getConnection() (*amqp.Connection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.conn != nil && !r.conn.IsClosed() {
		return r.conn, nil
	}

	if r.conn != nil {
		r.conn.Close()
	}

	conn, err := amqp.DialConfig(r.Cfg.Dns, amqp.Config{Vhost: r.Cfg.Vhost, Heartbeat: time.Duration(r.Cfg.Heartbeat) * time.Second})
	if err != nil {
		log.Printf("rabbitmq 创建连接失败 err=%v", err)
		return nil, err
	}
	r.conn = conn
	return conn, nil
}

// 创建一个新的操作对象
func NewRabbitMQ(cfg *RabbitMqConfig) *RabbitMQ {
	retryCnfs := make(map[int]int64, 0)
	if cfg.RetryCnfs != "" {
		cnfs := strings.Split(strings.Trim(cfg.RetryCnfs, ","), ",")
		for id, cnf := range cnfs {
			second, _ := strconv.ParseInt(cnf, 10, 64)
			retryCnfs[id] = second
		}
	}
	return &RabbitMQ{
		Cfg:       cfg,
		Consumers: make(map[*BusinessConfig]Handle),
		Lock:      sync.Mutex{},
		retryCnfs: retryCnfs,
		retryNum:  len(cfg.RetryCnfs),
		done:      make(chan struct{}),
	}
}

// Close 优雅关闭所有消费者协程和连接
func (r *RabbitMQ) Close() {
	close(r.done)

	r.mu.Lock()
	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}
	r.mu.Unlock()
}

func (r *RabbitMQ) Publish(b *BusinessConfig, body string) (err error) {
	conn, err := r.getConnection()
	if err != nil {
		log.Printf("rabbitmq Publish 链接获取失败 err=%v", err)
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("rabbitmq Publish 创建管道失败 err=%v", err)
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(b.GroupId, b.Topic, true, false, false, false, nil)
	if err != nil {
		log.Printf("发送消息创建交换机失败 err=%v", err)
		return err
	}

	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err = ch.QueueDeclare(b.Name, true, false, false, false, nil)
	if err != nil {
		log.Printf("发送消息创建队列失败 err=%v", err)
		return err
	}
	// 绑定任务
	err = ch.QueueBind(b.Name, b.Name, b.GroupId, false, nil)
	if err != nil {
		log.Printf("发送消息队列绑定交换机失败 err=%v", err)
		return err
	}

	err = ch.Publish("", b.Name, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(body),
		DeliveryMode: amqp.Persistent,
		Headers:      map[string]interface{}{"retry": 0},
	})
	if err != nil {
		return GeneralMessageDeliveryFailed
	}

	return nil
}

func (r *RabbitMQ) DeferPublish(b *BusinessConfig, body string, t time.Duration) (err error) {
	conn, err := r.getConnection()
	if err != nil {
		log.Printf("rabbitmq DeferPublish 链接获取失败 err=%v", err)
		return err
	}
	// 创建信道
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("rabbitmq DeferPublish 创建管道失败 err=%v", err)
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(b.GroupId, b.Topic, true, false, false, false, nil)
	if err != nil {
		log.Printf("发送延迟消息创建交换机失败 err=%v", err)
		return err
	}

	if t.Milliseconds() <= time.Second.Milliseconds() {
		log.Printf("发送延时消息，延迟时间参数是必须的填写")
		return errors.New("发送延时消息，延迟时间参数是必须的填写")
	}
	ttl := t.Milliseconds()
	delayQueueName := r.delayQueueName(b.Name, ttl)
	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err = ch.QueueDeclare(delayQueueName, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    b.GroupId,
		"x-message-ttl":             ttl, //消息存活时间
		"x-dead-letter-routing-key": delayQueueName,
	})
	if err != nil {
		log.Printf("发送延迟消息创建队列失败 err=%v", err)
		return err
	}
	// 绑定任务
	err = ch.QueueBind(b.Name, delayQueueName, b.GroupId, false, nil)
	if err != nil {
		log.Printf("发送延迟消息队列绑定交换机失败 err=%v", err)
		return err
	}

	err = ch.Publish("", delayQueueName, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(body),
		DeliveryMode: amqp.Persistent,
		Headers:      map[string]interface{}{"retry": 0},
	})
	if err != nil {
		log.Printf("发送延迟消息失败 err=%v", err)
		return DelayedMessageDeliveryFailed
	}
	return nil
}

func (r *RabbitMQ) delayQueueName(name string, ttl int64) string {
	var max int64
	var min int64
	// 0-15秒
	if ttl <= 15000 {
		min = 0
		max = 15000
	}
	// 15-30秒
	if ttl <= 30000 && ttl > 15000 {
		min = 15000
		max = 30000
	}
	// 30-5分钟
	if ttl <= 300000 && ttl > 30000 {
		min = 30000
		max = 300000
	}

	// 5-10分钟
	if ttl <= 600000 && ttl > 300000 {
		min = 300000
		max = 600000
	}

	// 10-15分钟
	if ttl <= 900000 && ttl > 600000 {
		min = 600000
		max = 900000
	}

	// 10-30分钟
	if ttl <= 1800000 && ttl > 900000 {
		min = 900000
		max = 1800000
	}

	// 30-60分钟
	if ttl <= 3600000 && ttl > 1800000 {
		min = 1800000
		max = 3600000
	}

	// 60-90分钟
	if ttl <= 5400000 && ttl > 3600000 {
		min = 3600000
		max = 5400000
	}

	// 1.5-3小时
	if ttl <= 10800000 && ttl > 5400000 {
		min = 1800000
		max = 5400000
	}
	// 3-7小时
	if ttl <= 25200000 && ttl > 10800000 {
		min = 10800000
		max = 25200000
	}
	// 7-10小时
	if ttl <= 36000000 && ttl > 25200000 {
		min = 25200000
		max = 36000000
	}
	// 10-15小时
	if ttl <= 54000000 && ttl > 36000000 {
		min = 36000000
		max = 54000000
	}
	// 15-24小时
	if ttl <= 86400000 && ttl > 54000000 {
		min = 54000000
		max = 86400000
	}

	// 24小时-3天
	if ttl <= 259200000 && ttl > 86400000 {
		min = 86400000
		max = 259200000
	}

	// 3天-7天
	if ttl <= 604800000 && ttl > 259200000 {
		min = 259200000
		max = 604800000
	}
	// 7天-14天
	if ttl <= 1209600000 && ttl > 604800000 {
		min = 604800000
		max = 1209600000
	}

	// 14天-20天
	if ttl <= 1728000000 && ttl > 1209600000 {
		min = 1209600000
		max = 1728000000
	}
	// 20天-30天
	if ttl <= 2592000000 && ttl > 1728000000 {
		min = 1728000000
		max = 2592000000
	}
	// 30天-90天
	if ttl <= 7776000000 && ttl > 2592000000 {
		min = 2592000000
		max = 7776000000
	}

	// 90天-180天
	if ttl <= 15552000000 && ttl > 7776000000 {
		min = 7776000000
		max = 15552000000
	}
	// 180天-360天
	if ttl <= 31104000000 && ttl > 15552000000 {
		min = 15552000000
		max = 31104000000
	}
	// 360天-无限期
	if ttl > 31104000000 {
		min = 31104000000
		max = -1
	}

	return fmt.Sprintf("enqueue.%s.%d.%d.x.delay", name, min, max)
}

func (r *RabbitMQ) Listen() {
	if len(r.Consumers) <= 0 {
		log.Printf("rabbitmq消费者监听成功,共%d个消费者", len(r.Consumers))
		return
	}

	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Printf("rabbitmq listener panic err=%v", panicErr)
		}
	}()

	var wg sync.WaitGroup
	for business, handle := range r.Consumers {
		b := *business
		h := handle
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.consumeWithReconnect(&b, h)
		}()
	}

	// 阻塞等待退出信号
	<-r.done
	log.Printf("rabbitmq 收到退出信号，等待消费者协程退出...")
	wg.Wait()
	log.Printf("rabbitmq 所有消费者协程已退出")
}

func (r *RabbitMQ) Register(b *BusinessConfig, handle Handle) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	r.Consumers[b] = handle
}

// consumeWithReconnect 带重连机制的消费者，从单例连接创建独立channel
func (r *RabbitMQ) consumeWithReconnect(b *BusinessConfig, handle Handle) {
	for {
		// 检查是否收到退出信号
		select {
		case <-r.done:
			return // 收到信号，直接退出
		default:
		}
		// 获取mq 连接
		conn, err := r.getConnection()
		if err != nil {
			log.Printf("消费者[%s]获取连接失败，3秒后重试 err=%v", b.Name, err)
			r.waitOrDone(3 * time.Second)
			continue
		}

		r.receiver(conn, b, handle)

		select {
		case <-r.done:
			return
		default:
			log.Printf("消费者[%s]channel断开，3秒后重连", b.Name)
			r.waitOrDone(3 * time.Second)
		}
	}
}

// receiver 消费消息，从共享连接创建独立channel
func (r *RabbitMQ) receiver(conn *amqp.Connection, b *BusinessConfig, handle Handle) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("消费者[%s]创建信道失败 err=%v", b.Name, err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(b.GroupId, b.Topic, true, false, false, false, nil)
	if err != nil {
		log.Printf("消费者[%s]创建交换机失败 err=%v", b.Name, err)
		return
	}

	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err = ch.QueueDeclare(b.Name, true, false, false, false, nil)
	if err != nil {
		log.Printf("消费者[%s]创建队列失败 err=%v", b.Name, err)
		return
	}
	// 绑定任务
	err = ch.QueueBind(b.Name, b.Name, b.GroupId, false, nil)
	if err != nil {
		log.Printf("消费者[%s]队列绑定交换机失败 err=%v", b.Name, err)
		return
	}

	msgList, err := ch.Consume(b.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("消费者[%s]消费消息失败 err=%v", b.Name, err)
		return
	}

	log.Printf("消费者[%s]开始消费", b.Name)
	for msg := range msgList {
		// 处理数据
		err := handle(string(msg.Body))
		if err == nil {
			// 确认消息,必须为false   false为单条确认，true批量确认（可能会丢消息）
			err = msg.Ack(false)
			if err != nil {
				log.Printf("消息消费ack失败 err=%v", err)
			}
		}
		if err != nil {
			tempRetry, ok := msg.Headers["retry"]
			var retry int = 0
			if ok {
				retry, _ = strconv.Atoi(fmt.Sprintf("%d", tempRetry))
				retry = retry + 1
			}

			//消息处理失败 进入延时尝试机制
			if r.retryNum != 0 && retry <= r.retryNum && ok {
				if len(r.retryCnfs) > 0 && retry > 0 && retry-1 < len(r.retryCnfs) {
					retrySecond := r.retryCnfs[retry-1]
					r.publishRetry(string(msg.Body), retry, retrySecond, b)
				}

			} else {
				//消息失败 入库db
				log.Printf("消息处理多次后还是失败了 可以扩展功能写入到 db 或邮件通知")
				//receiver.FailAction(err, msg.Body)
			}
			// 确认消息,必须为false   false为单条确认，true批量确认（可能会丢消息）
			err = msg.Ack(false)
			if err != nil {
				log.Printf("消息消费ack失败 err=%v", err)
			}
		}
	}
}

// waitOrDone 等待指定时间，或收到退出信号时提前返回
func (r *RabbitMQ) waitOrDone(d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-r.done:
	case <-timer.C:
	}
}

// 发送重试消息（复用单例连接）
func (r *RabbitMQ) publishRetry(body string, retry int, retrySecond int64, b *BusinessConfig) {
	conn, err := r.getConnection()
	if err != nil {
		log.Printf("发送重试消息连接失败 err=%v", err)
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("发送重试消息创建信道失败 err=%v", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(b.GroupId, b.Topic, true, false, false, false, nil)
	if err != nil {
		log.Printf("发送重试消息创建交换机失败 err=%v", err)
		return
	}
	ttl := retrySecond * 1000
	retryQueueName := fmt.Sprintf("enqueue.%s.%v.x.retry", b.Name, ttl)

	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err = ch.QueueDeclare(retryQueueName, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    b.GroupId,
		"x-message-ttl":             ttl, //消息存活时间
		"x-dead-letter-routing-key": retryQueueName,
	})
	if err != nil {
		log.Printf("发送重试消息创建队列失败 err=%v", err)
		return
	}
	// 绑定任务
	err = ch.QueueBind(b.Name, retryQueueName, b.GroupId, false, nil)
	if err != nil {
		log.Printf("发送重试消息队列绑定交换机失败 err=%v", err)
		return
	}
	err = ch.Publish("", retryQueueName, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         []byte(body),
		Headers:      map[string]interface{}{"retry": retry},
	})
	if err != nil {
		log.Printf("发送重试消息失败 err=%v", err)
		return
	}
}
