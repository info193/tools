package pay

// 不同类型pay的实现接口
type IPAY interface {
	Wechat() *WechatSDK
	Alipay() *AlipaySDK
}

func NewPay(cfg *PayConfig) IPAY {
	switch cfg.PayEngine {
	case "self":
		return NewSelf(cfg)
	default:
		panic("New Mq Error")
	}
	return nil
}
