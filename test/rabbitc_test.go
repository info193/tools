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
		Dns:       "amqp://rabbit_prod:UsUngiYtaGG5QMqK@192.168.7.73:5672/",
		Vhost:     "/",
		Heartbeat: 5,
		RetryCnfs: []int64{10, 30, 60},
	}
	rabbit := mq.NewRabbitMQ(conf)

	//var cns CNss
	b := mq.NewBusiness("demo", "direct", "demo", "demo")
	rabbit.Register(b, Coner)
	rabbit.Listen()

}
