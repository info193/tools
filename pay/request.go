package pay

import (
	"github.com/go-pay/gopay/wechat/v3"
	"time"
)

// PayRequest 支付请求
type PayRequest struct {
	Subject    string        `json:"subject"`      // 订单标题 商品描述
	Amount     int64         `json:"amount"`       // 付款金额（分）
	OutTradeNo string        `json:"out_trade_no"` // 订单号
	TimeExpire time.Duration `json:"time_expire"`  // 订单过期时间  秒
	NotifyUrl  string        `json:"notify_url"`   // 异步回调地址
	RefundUrl  string        `json:"refund_url"`   // 同步回调地址
	Attach     string        `json:"attach"`       // 扩展字段
	Account    string        `json:"account"`      // 支付账号
}

// QueryRequest 交易查询请求
type QueryRequest struct {
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
	TradeNo    string `json:"trade_no"`     // 原支付单号(支付平台)
}

// RefundQueryRequest 退款查询请求
type RefundQueryRequest struct {
	OutRefundNo string `json:"out_refund_no"` // 退款请求号,如果在退款请求时未传入，则该值为创建交易时的商户订单号。
	OutTradeNo  string `json:"out_trade_no"`  // 商户订单号
	TradeNo     string `json:"trade_no"`      // 原支付单号(支付平台)
}

type Response struct {
	Mode          int64                `json:"mode"` // 模式 1自有
	AppWechatInfo *wechat.AppPayParams `json:"app_wechat_info"`
	AppAlipayInfo string               `json:"app_alipay_info"`
	AppletWechat  *wechat.AppletParams `json:"applet_wechat"`
}
