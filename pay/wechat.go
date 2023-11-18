package pay

import (
	"context"
	"errors"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"os"
	"time"
)

// 定义WechatSDK对象
type WechatSDK struct {
	cfg  *Wechat
	conn *wechat.ClientV3
}

func (l *WechatSDK) client() error {
	content, err := os.ReadFile(l.cfg.MchSecretCert)
	if err != nil {
		return CertificateFail
	}

	l.conn, err = wechat.NewClientV3(l.cfg.MchId, l.cfg.SerialNo, l.cfg.MchSecretKey, string(content))
	if err != nil {
		return err
	}
	// 启用自动同步返回验签，并定时更新微信平台API证书（开启自动验签时，无需单独设置微信平台API证书和序列号）
	err = l.conn.AutoVerifySign()
	if err != nil {
		return AutosignatureFail
	}
	if l.cfg.DebugSwitch == 1 {
		l.conn.DebugSwitch = gopay.DebugOn
	}
	if l.cfg.DebugSwitch == 0 {
		l.conn.DebugSwitch = gopay.DebugOff
	}
	return nil
}

// App 支付
func (l *WechatSDK) App(request *PayRequest) (*Response, error) {
	if l.conn == nil {
		return nil, PayGatewayFail
	}
	expire := time.Now().Add(request.TimeExpire * time.Second).Format(time.RFC3339)
	var notifyUrl string
	if request.NotifyUrl != "" {
		notifyUrl = l.cfg.NotifyUrl
	} else {
		notifyUrl = request.NotifyUrl
	}
	bm := make(gopay.BodyMap)
	bm.Set("appid", l.cfg.Appid).
		Set("description", request.Subject).
		Set("out_trade_no", request.OutTradeNo).
		Set("time_expire", expire).
		Set("notify_url", notifyUrl).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", request.Amount).Set("currency", "CNY")
		})

	wxRsp, err := l.conn.V3TransactionApp(context.Background(), bm)
	if err != nil {
		return nil, err
	}
	app, err := l.conn.PaySignOfApp(l.cfg.Appid, wxRsp.Response.PrepayId)
	if err != nil {
		return nil, err
	}

	return &Response{
		Mode:          1,
		AppWechatInfo: app,
	}, nil
}

// Mini 微信小程序支付
func (l *WechatSDK) Mini(request *PayRequest) (*Response, error) {
	if l.conn == nil {
		return nil, PayGatewayFail
	}
	expire := time.Now().Add(request.TimeExpire * time.Second).Format(time.RFC3339)
	var notifyUrl string
	if request.NotifyUrl != "" {
		notifyUrl = l.cfg.NotifyUrl
	} else {
		notifyUrl = request.NotifyUrl
	}
	bm := make(gopay.BodyMap)
	bm.Set("appid", l.cfg.MiniAppid).
		Set("description", request.Subject).
		Set("out_trade_no", request.OutTradeNo).
		Set("time_expire", expire).
		Set("notify_url", notifyUrl).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", request.Amount).Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", request.Account)
		})

	resp, err := l.conn.V3TransactionJsapi(context.Background(), bm)
	if err != nil {
		return nil, err

	}
	applet, err := l.conn.PaySignOfApplet(l.cfg.MiniAppid, resp.Response.PrepayId)
	if err != nil {
		return nil, err
	}

	return &Response{
		Mode:         1,
		AppletWechat: applet,
	}, nil
}

// Query 交易查询
func (l *WechatSDK) Query(request *QueryRequest) (*wechat.QueryOrder, error) {
	if l.conn == nil {
		return nil, PayGatewayFail
	}
	if request.TradeNo == "" && request.OutTradeNo == "" {
		return nil, ParamFail
	}

	var orderNo string
	var orderType wechat.OrderNoType
	if request.TradeNo != "" {
		orderNo = request.TradeNo
		orderType = wechat.TransactionId
	}
	if request.OutTradeNo != "" {
		orderNo = request.OutTradeNo
		orderType = wechat.OutTradeNo
	}
	response, err := l.conn.V3TransactionQueryOrder(context.Background(), orderType, orderNo)
	if err != nil {
		return nil, err
	}
	// 判断请求状态
	if response.Code != wechat.Success {
		return nil, errors.New(response.Error)
	}

	return response.Response, nil
}

// RefundQuery 退款订单查询
func (l *WechatSDK) RefundQuery(request *RefundQueryRequest) (*wechat.RefundQueryResponse, error) {
	if l.conn == nil {
		return nil, PayGatewayFail
	}
	if request.OutRefundNo == "" {
		return nil, ParamFail
	}

	response, err := l.conn.V3RefundQuery(context.Background(), request.OutRefundNo, nil)
	if err != nil {
		return nil, err
	}
	// 判断请求状态
	if response.Code != wechat.Success {
		return nil, errors.New(response.Error)
	}

	return response.Response, nil
}
