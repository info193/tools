package test

import (
	"fmt"
	"github.com/info193/tools/mq"
	"testing"
)

func TestRabbitp(t *testing.T) {
	//
	//conf := &mq.RabbitMqConfig{
	//	HostSource: "amqp://rabbit_prod:UsUngiYtaGG5QMqK@192.168.7.73:5672/",
	//	Vhost:      "/",
	//	Heartbeat:  5,
	//	RetryCnf:   map[int64]int64{1: 0, 2: 30, 3: 60},
	//}
	//rc := make([]mq.RetryCnf, 0)
	//rc = append(rc, mq.RetryCnf{
	//	Num: 1,
	//	Ttl: 10,
	//})
	//rc = append(rc, mq.RetryCnf{
	//	Num: 2,
	//	Ttl: 30,
	//})
	//rc = append(rc, mq.RetryCnf{
	//	Num: 3,
	//	Ttl: 60,
	//})

	conf := &mq.RabbitMqConfig{
		Dns:       "amqp://rabbit_prod:UsUngiYtaGG5QMqK@192.168.7.73:5672/",
		Vhost:     "/",
		Heartbeat: 5,
		RetryCnfs: []int64{10, 30, 60},
	}
	rabbit := mq.NewRabbitMQ(conf)
	//b := mq.NewBusiness("test", "direct", "test", "test")
	b := mq.NewBusiness("demo", "direct", "demo", "demo")

	err := rabbit.Publish(b, "cesiyixia")
	fmt.Println(err, "++++++++++++++++")
	//err = rabbit.DeferPublish(b, "yanchixiaoxi", time.Second*15)
	fmt.Println(err, "-------------")

}
