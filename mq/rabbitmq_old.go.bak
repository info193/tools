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
	Cfg        *RabbitMqConfig // 配置
	connection *amqp.Connection
	Channel    *amqp.Channel
	Consumers  map[*BusinessConfig]Handle // 消费者
	Lock       sync.Mutex                 // 锁
	retryCnfs  map[int]int64
	retryNum   int
}

// 链接rabbitMQ
func (r *RabbitMQ) MqConnect() (err error) {
	mqConn, err := amqp.DialConfig(r.Cfg.Dns, amqp.Config{Vhost: r.Cfg.Vhost, Heartbeat: time.Duration(r.Cfg.Heartbeat) * time.Second})
	//mqConn, err = amqp.Dial(r.dns)
	r.connection = mqConn // 赋值给RabbitMQ对象
	if err != nil {
		log.Printf("rabbitmq 创建mq链接失败 err=%v", err)
	}
	return
}

// 关闭mq链接
func (r *RabbitMQ) CloseMqConnect() (err error) {
	err = r.connection.Close()
	if err != nil {
		log.Printf("关闭mq链接失败 err=%v", err)
	}
	return
}

// 链接rabbitMQ
func (r *RabbitMQ) MqOpenChannel() (err error) {
	mqConn := r.connection
	r.Channel, err = mqConn.Channel()
	if err != nil {
		log.Printf("rabbitmq创建管道失败 err=%v", err)
	}
	return err
}

// 链接rabbitMQ
func (r *RabbitMQ) CloseMqChannel() (err error) {
	err = r.Channel.Close()
	if err != nil {
		log.Printf("关闭rabbitmq链接失败 err=%v", err)
	}
	return err
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
	}
}

func (r *RabbitMQ) Publish(b *BusinessConfig, body string) (err error) {
	// 1.创建mq链接
	err = r.MqConnect()
	if err != nil {
		return err
	}
	defer func() {
		_ = r.CloseMqConnect()
	}()
	// 2.创建信道
	err = r.MqOpenChannel()
	ch := r.Channel
	if err != nil {
		log.Printf("发送消息创建channel信道失败 err=%v", err)
		return err
	}
	defer func() {
		_ = r.Channel.Close()
	}()

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

	err = r.Channel.Publish("", b.Name, false, false, amqp.Publishing{
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
	// 创建mq链接
	err = r.MqConnect()
	if err != nil {
		return err
	}
	defer func() {
		_ = r.CloseMqConnect()
	}()

	// 创建信道
	err = r.MqOpenChannel()
	ch := r.Channel
	if err != nil {
		log.Printf("发送延迟消息创建信道失败 err=%v", err)
		return err
	}
	defer r.Channel.Close()

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
	//delayQueueName := fmt.Sprintf("enqueue.%s.%v.x.delay", b.Name, ttl)
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

	err = r.Channel.Publish("", delayQueueName, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(body),
		DeliveryMode: amqp.Persistent, //  将消息标记为持久消息
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
	if ttl <= 15552000000 && ttl > 7776000000 {
		min = 7776000000
		max = 15552000000
	}
	// 360天-无限期
	if ttl > 15552000000 {
		min = 15552000000
		max = -1
	}

	return fmt.Sprintf("enqueue.%s.%d.%d.x.delay", name, min, max)
}

func (r *RabbitMQ) Listen() {
	var (
		exitTask bool
	)
	if len(r.Consumers) <= 0 {
		log.Printf("rabbitmq消费者监听成功,共%d个消费者", len(r.Consumers))
		return
	}

	//链接rabbitMQ
	err := r.MqConnect()
	if err != nil {
		return
	}

	defer func() {
		if panicErr := recover(); panicErr != nil {
			fmt.Println(recover())
			err = errors.New(fmt.Sprintf("%s", panicErr))
		}
	}()

	//开始执行任务
	for business, handle := range r.Consumers {
		b := *business
		h := handle
		go r.do(&b, h)
	}

	scheduleTimer := time.NewTimer(time.Millisecond * 300)
	exitTask = true
	for {
		select {
		case <-scheduleTimer.C:
			//fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~")
		}
		// 重置调度间隔
		scheduleTimer.Reset(time.Millisecond * 300)
		if !exitTask {
			break
		}
	}
	fmt.Println("exit")
	return
}
func (r *RabbitMQ) Register(b *BusinessConfig, handle Handle) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	r.Consumers[b] = handle
}

// 监听接收者接收任务
func (r *RabbitMQ) Receiver(b *BusinessConfig, handle Handle) {
	err := r.MqOpenChannel()
	ch := r.Channel
	if err != nil {
		log.Printf("消费者创建信道失败 err=%v", err)
		return
	}
	defer r.Channel.Close()
	err = ch.ExchangeDeclare(b.GroupId, b.Topic, true, false, false, false, nil)
	if err != nil {
		log.Printf("消费者创建交换机失败 err=%v", err)
		return
	}

	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err = ch.QueueDeclare(b.Name, true, false, false, false, nil)
	if err != nil {
		log.Printf("消费者创建队列失败 err=%v", err)
		return
	}
	// 绑定任务
	err = ch.QueueBind(b.Name, b.Name, b.GroupId, false, nil)
	if err != nil {
		log.Printf("消费者队列绑定交换机失败 err=%v", err)
		return
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	//err = ch.Qos(1, 0, false)
	msgList, err := ch.Consume(b.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("消费者消费消息失败 err=%v", err)
		return
	}
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

func (r *RabbitMQ) do(b *BusinessConfig, receiver Handle) {
	// 验证链接是否正常
	err := r.MqOpenChannel()
	if err != nil {
		return
	}
	r.Receiver(b, receiver)
}

// 发送重试消息
func (r *RabbitMQ) publishRetry(body string, retry int, retrySecond int64, b *BusinessConfig) {
	_ = r.MqConnect()
	defer func() {
		_ = r.CloseMqConnect()
	}()

	err := r.MqOpenChannel()
	ch := r.Channel
	if err != nil {
		log.Printf("发送重试消息创建信道失败 err=%v", err)
		return
	}
	defer r.Channel.Close()

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
	err = r.Channel.Publish("", retryQueueName, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent, //  将消息标记为持久消息
		Body:         []byte(body),
		Headers:      map[string]interface{}{"retry": retry},
	})
	if err != nil {
		log.Printf("发送重试消息失败 err=%v", err)
		return
	}

}
