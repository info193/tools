package pay

import (
	"context"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/shopspring/decimal"
	"os"
	"time"
)

// 定义AlipaySDK对象
type AlipaySDK struct {
	cfg  *Alipay
	conn *alipay.Client
}

func (l *AlipaySDK) client() error {
	content, err := os.ReadFile(l.cfg.AlipayPublicCertPath)
	if err != nil {
		return CertificateFail
	}

	l.conn, err = alipay.NewClient(l.cfg.Appid, l.cfg.SecretCert, l.cfg.IsProd)
	if err != nil {
		return err
	}
	// 打开Debug开关，输出日志，默认关闭
	if l.cfg.DebugSwitch == 1 {
		l.conn.DebugSwitch = gopay.DebugOn
	}
	if l.cfg.DebugSwitch == 0 {
		l.conn.DebugSwitch = gopay.DebugOff
	}

	l.conn.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
							SetCharset(alipay.UTF8).       // 设置字符编码，不设置默认 utf-8
							SetSignType(l.cfg.SignType).   // 设置签名类型，不设置默认 RSA2
							SetReturnUrl(l.cfg.ReturnUrl). // 设置返回URL
							SetNotifyUrl(l.cfg.NotifyUrl)  // 设置异步通知URL
	//SetAppAuthToken()              // 设置第三方应用授权
	// 自动同步验签（只支持证书模式）
	// 传入 alipayPublicCert.crt 内容
	l.conn.AutoVerifySign(content)
	// 公钥证书模式，需要传入证书，以下两种方式二选一
	// 证书路径
	err = l.conn.SetCertSnByPath(l.cfg.AppPublicCertPath, l.cfg.AlipayRootCertPath, l.cfg.AlipayPublicCertPath)
	// 证书内容
	//err := client.SetCertSnByContent("appPublicCert.crt bytes", "alipayRootCert bytes", "alipayPublicCert.crt bytes")
	if err != nil {
		return err
	}

	return nil
}

// app 支付
func (l *AlipaySDK) App(request *PayRequest) (*Response, error) {
	if l.conn == nil {
		return nil, PayGatewayFail
	}
	tempAmount := decimal.NewFromInt(request.Amount)
	amount := fmt.Sprintf("%.2f", tempAmount.Div(decimal.NewFromInt(100)).InexactFloat64())
	timeoutExpress := time.Minute * 5
	if request.TimeExpire > 0 {
		timeoutExpress = request.TimeExpire
	}

	bm := make(gopay.BodyMap)
	bm.Set("subject", request.Subject).
		Set("out_trade_no", request.OutTradeNo).
		Set("total_amount", amount).
		Set("timeout_express", timeoutExpress)
	response, err := l.conn.TradeAppPay(context.Background(), bm)
	if err != nil {
		if bizErr, ok := alipay.IsBizError(err); ok {
			return nil, bizErr
		}
	}

	return &Response{
		Mode:          1,
		AppAlipayInfo: response,
	}, nil
}

// 交易订单查询
func (l *AlipaySDK) Query(request *QueryRequest) (*alipay.TradeQuery, error) {
	if l.conn == nil {
		return nil, PayGatewayFail
	}
	if request.TradeNo == "" && request.OutTradeNo == "" {
		return nil, ParamFail
	}

	bm := make(gopay.BodyMap)
	if request.TradeNo != "" {
		bm.Set("trade_no", request.TradeNo)
	}
	if request.OutTradeNo != "" {
		bm.Set("out_trade_no", request.OutTradeNo)
	}
	response, err := l.conn.TradeQuery(context.Background(), bm)
	if err != nil {
		if bizErr, ok := alipay.IsBizError(err); ok {
			return nil, bizErr
		}
	}
	return response.Response, nil
}

// 退款订单查询
func (l *AlipaySDK) RefundQuery(request *RefundQueryRequest) (*alipay.TradeRefundQuery, error) {
	if l.conn == nil {
		return nil, PayGatewayFail
	}
	if (request.TradeNo == "" && request.OutTradeNo == "") || request.OutRefundNo == "" {
		return nil, ParamFail
	}

	bm := make(gopay.BodyMap)
	bm.Set("out_request_no", request.OutRefundNo)
	if request.TradeNo != "" {
		bm.Set("trade_no", request.TradeNo)
	}
	if request.OutTradeNo != "" {
		bm.Set("out_trade_no", request.OutTradeNo)
	}
	response, err := l.conn.TradeFastPayRefundQuery(context.Background(), bm)
	if err != nil {
		if bizErr, ok := alipay.IsBizError(err); ok {
			return nil, bizErr
		}
	}
	return response.Response, nil
}
