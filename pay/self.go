package pay

type Self struct {
	cfg *PayConfig
}

func NewSelf(c *PayConfig) *Self {
	return &Self{cfg: c}
}
func (l *Self) Wechat() *WechatSDK {
	if wechat, ok := l.cfg.Wechat[l.cfg.WechatMode]; ok {
		wechatSDK := &WechatSDK{cfg: &wechat}
		wechatSDK.client()
		return wechatSDK
	}
	return nil
}
func (l *Self) Alipay() *AlipaySDK {
	if alipay, ok := l.cfg.Alipay[l.cfg.AlipayMode]; ok {
		alipaySDK := &AlipaySDK{cfg: &alipay}
		alipaySDK.client()
		return alipaySDK
	}
	return nil
}
