package mq

type RocketMqAliConfig struct {
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Endpoint   string `json:"endpoint"` //设置HTTP协议客户端接入点，进入消息队列RocketMQ版控制台实例详情页面的接入点区域查看。
	InstanceId string `json:"instance_id"`
}

type RocketMqConfig struct {
	Endpoint string `json:"endpoint,omitempty"`
}

type NsqConfig struct {
	Lookupds []string `json:"lookupds,omitempty"`
}

type RabbitMqConfig struct {
	HostSource string `json:"host_source"`
	Vhost      string `json:"vhost"`
	//Retry      int64           `json:"retry"`
	Heartbeat int64         `json:"heartbeat"`
	RetryCnf  map[int]int64 `json:"retry_delay"`
}
