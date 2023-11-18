package pay

import "errors"

var (
	PayGatewayFail    = errors.New("支付网关链接错误")
	CertificateFail   = errors.New("证书错误")
	AutosignatureFail = errors.New("自动验签失败")
	ParamFail         = errors.New("请求参数有误")
)
