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

	//str := "0,20,30"
	//if str != "" {
	//	arr := strings.Split(str, ",")
	//	//array := make([]int64, 0)
	//	retryCnfs := make(map[int]int64, 0)
	//	for id, value := range arr {
	//		num, _ := strconv.ParseInt(value, 10, 64)
	//		retryCnfs[id] = num
	//	}
	//	fmt.Println(retryCnfs, len(retryCnfs))
	//}
	//return

	conf := &mq.RabbitMqConfig{
		Dns:       "amqp://rabbit_prod:UsUngiYtaGG5QMqK@192.168.7.73:5672/",
		Vhost:     "/",
		Heartbeat: 5,
		RetryCnfs: "0,20,30", // 第一次重试 0秒，第二次重试 20秒，第三次重试30秒
	}
	rabbit := mq.NewRabbitMQ(conf)
	//b := mq.NewBusiness("test", "direct", "test", "test")
	b := mq.NewBusiness("demo", "direct", "demo", "demo")

	err := rabbit.Publish(b, "cesiyixia")
	fmt.Println(err, "++++++++++++++++")
	//err = rabbit.DeferPublish(b, "yanchixiaoxi", time.Second*15)
	fmt.Println(err, "-------------")

}
