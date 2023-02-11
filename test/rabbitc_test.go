package test

import (
	errs "errors"
	"fmt"
	"github.com/info193/tools/mq"
	"testing"
)

type CNss struct {
}

func (c CNss) Consumer(body []byte) error {

	fmt.Println("Consumer--", string(body))
	return errs.New("消费错误 请重试.")
}
func (c CNss) FailAction(err error, body []byte) error {
	fmt.Println("FailAction", err, string(body))
	return nil
}

func Coner(body string) error {
	fmt.Println("Consumer--", body)
	return errs.New("消费错误 请重试.")
}

func TestRabbitc(t *testing.T) {
	//
	//conf := &mq.RabbitMqConfig{
	//	HostSource: "amqp://rabbit_prod:UsUngiYtaGG5QMqK@192.168.7.73:5672/",
	//	Vhost:      "/",
	//	Heartbeat:  5,
	//	RetryCnf:   map[int64]int64{1: 0, 2: 30, 3: 60},
	//}

	conf := &mq.RabbitMqConfig{
		HostSource: "amqp://rabbit_prod:UsUngiYtaGG5QMqK@192.168.7.73:5672/",
		Vhost:      "/",
		Heartbeat:  5,
		RetryCnf:   map[int]int64{1: 10, 2: 30, 3: 60},
	}
	rabbit := mq.NewRabbit(conf)

	//var cns CNss
	b := mq.NewBusiness("demo", "direct", "demo", "demo")
	rabbit.Register(b, Coner)
	rabbit.Listen()
	//rabbitmq := rabbit.Recvs(b, cns)
	//fmt.Println(rabbitmq)

	//err := rabbit.Sends(b, "cesiyixia")
	//fmt.Println(err, "-------------")
	//qe := mq.QueueExchange{
	//	QuName: "demo",
	//	RtKey:  "demo",
	//	ExName: "demo",
	//	ExType: "direct",
	//	Dns:    "amqp://rabbit_prod:UsUngiYtaGG5QMqK@192.168.7.73:5672/",
	//}
	//
	//rabbitmq := mq.Send(qe, "cesiyixia")
	////mq.Send(qe, "cesiyixia2222")
	//mq.SendDelay(qe, "55555", 15)
	//fmt.Println(rabbitmq)
	//defer func() {
	//	//rabbitmq.CloseMqChannel()
	//	rabbitmq.CloseMqConnect()
	//}()
	//
	//b := mq.NewBusiness("demo", "direct", "demo", "demo")
	////rabbitmq.Publish(b, fmt.Sprintf("demo iiwiiiii%v", time.Now().Format("2006-01-02 15:04:05")))
	//rabbitmq.Publish(b, "demo")
	////rabbitmq.Publish(b, "test1")
	////rabbitmq.Publish(b, "test2")
	//
	//bss := mq.NewBusiness("test", "direct", "test", "test")
	//rabbitmq.Publish(bss, "test")
	//rabbitmq.Publish(bss, "test3")

	//bs := mq.NewBusiness("demo", "direct", "demo", "demo")
	//rabbitmq.Publish(b, fmt.Sprintf("demo iiwiiiii%v", time.Now().Format("2006-01-02 15:04:05")))

	//defer rabbitmq.Close()

	//rabbitmq.Publish(b, "test2")
	//rabbitmq.Publish(b, "test3")

	//rabbitmq.DeferPublish(b, fmt.Sprintf("demo deoooooo%v", time.Now().Format("2006-01-02 15:04:05")), time.Second*10)

}
