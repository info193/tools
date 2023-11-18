package pay

// PayConfig 不同业务的配置
type PayConfig struct {
	PayEngine  string            `json:"pay_engine"` // 支付引擎self
	AlipayMode string            `json:"alipay_mode"`
	WechatMode string            `json:"wechat_mode"`
	Alipay     map[string]Alipay `json:"alipay"`
	Wechat     map[string]Wechat `json:"wechat"`
}

type Alipay struct {
	Appid                string `json:"appid"`                   // app_id
	SecretCert           string `json:"secret_cert"`             // 应用私钥 支持PKCS1和PKCS8
	IsProd               bool   `json:"is_prod"`                 // isProd：是否是正式环境，沙箱环境请选择新版沙箱应用。
	SignType             string `json:"sign_type"`               // 设置签名类型，不设置默认 RSA2
	ReturnUrl            string `json:"return_url"`              // 设置返回URL
	NotifyUrl            string `json:"notify_url"`              // 设置异步通知URL
	AppPublicCertPath    string `json:"app_public_cert_path"`    // 应用公钥证书 路径
	AlipayRootCertPath   string `json:"alipay_root_cert_path"`   // 支付宝根证书 路径
	AlipayPublicCertPath string `json:"alipay_public_cert_path"` // 支付宝公钥证书 路径
	DebugSwitch          int64  `json:"debug_switch"`            // 是否开启debug打印日志 0否 1是
}

type Wechat struct {
	Appid         string `json:"appid"`           // app_id 选填-app 的 appid
	MiniAppid     string `json:"mini_appid"`      // mini_app_id 选填-小程序 的 appid
	MchId         string `json:"mch_id"`          // app_id
	SerialNo      string `json:"serial_no"`       // 商户API证书的证书序列号
	MchSecretKey  string `json:"mch_secret_key"`  // 商户秘钥
	MchSecretCert string `json:"mch_secret_cert"` // 商户私钥 字符串或路径
	MchPublicCert string `json:"mch_public_cert"` // 商户公钥证书 字符串或路径
	DebugSwitch   int64  `json:"debug_switch"`    // 是否开启debug打印日志 0否 1是
	NotifyUrl     string `json:"notify_url"`      // 设置异步通知URL
}
