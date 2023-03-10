package mq

import (
	"time"
)

// Handle 消费者业务方法
type Handle func(string) error

// 不同类型mq的实现接口
type IMQ interface {
	// Publish 生产普通消息
	Publish(b *BusinessConfig, msg string) error
	// DeferPublish 生产延时消息
	DeferPublish(b *BusinessConfig, msg string, t time.Duration) error
	// Register 注册一个消费者
	Register(b *BusinessConfig, handle Handle)
	// Listen 消费者监听
	Listen()
}

// MqConfig
// @Description: 消息队列
type MqConfig struct {
	Switch bool
	Use    string
	//AsynqConfig       *AsynqConfig
	//RocketMqAliConfig *RocketMqAliConfig
	//RocketMqConfig    *RocketMqConfig
	//NsqConfig         *NsqConfig
	RabbitMqConfig *RabbitMqConfig
}

func NewMq(cfg *MqConfig) IMQ {
	switch cfg.Use {
	case "RabbitMQ":
		return NewRabbitMQ(cfg.RabbitMqConfig)
	case "RocketMQAli":
		//return NewRocketAli(cfg.RocketMqAliConfig)
	case "RocketMQ":
		//return NewRocket(cfg.RocketMqConfig)
	case "Nsq":
		//return NewNSQ(cfg.NsqConfig)
	case "Asynq":
		//return NewAsynq(cfg.AsynqConfig)
	default:
		panic("New Mq Error")
	}
	return nil
}
