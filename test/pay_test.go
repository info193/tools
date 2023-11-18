package test

import (
	"fmt"
	"github.com/info193/tools/pay"
	"testing"
)

func TestPay(t *testing.T) {
	alipay := make(map[string]pay.Alipay)
	wechat := make(map[string]pay.Wechat)
	alipay["default"] = pay.Alipay{
		Appid:                "2019013155890001",
		SecretCert:           "MIIEArO0wf8RmhjowQ/1igAWPcnoxOQxQgT94E0gbB9K7h2f4PY0guoX8VRBn7qKDqNUYihexsLF0lg93IvBOq2IIBAAKCAQEd/lrzl9gIwZNlSJ82RfoKMitLcbKE+M6IQ5v3mv3Kj0twx/Fm10gbB9Kh2f4PY0guoX8VRBnqKDqNUYihexsLFlg93I5vBOq2IIBAAKCAQEd/lrzl9gIwZNlS3J2RfoKMitLcbK/0gbB9Kh2f4P50guoX28VRBnqKDqNUYihexsLFlg93IvBdOq2IIBAAKCAQEd/lr24zl9gIwZNlSJ2RfoKMitLcbKs0gbB9Kh2f4PY0guoX8VRBnqKDqNUYihexsL3Flg93IvBOq2IIBAAKCAQEd/lrzl9gIwZ3NlSJ2RfoKMitLcbK",
		IsProd:               true,
		SignType:             "RSA2",
		ReturnUrl:            "https://www.xxx.com",
		NotifyUrl:            "https://www.xxx.com",
		AppPublicCertPath:    "../cert/alipay/appCertPublicKey_2019013155890001.crt",
		AlipayRootCertPath:   "../cert/alipay/alipayRootCert.crt",
		AlipayPublicCertPath: "../cert/alipay/alipayCertPublicKey_RSA2.crt",
		DebugSwitch:          0,
	}

	wechat["default"] = pay.Wechat{
		Appid:         "wx5c1050fk9p652612",
		MiniAppid:     "wx5afhd05wtbe1f28", // mini_app_id 选填-小程序 的 appid
		MchId:         "1551725412",
		SerialNo:      "6A0CC35F95D5301ABC29A81DB38A7149CA6384FE",
		MchSecretKey:  "Itc4V5udGXIpLH1JnUzT54ZswC1Sbjlq2",
		MchSecretCert: "../cert/apiclient_key.pem",
		MchPublicCert: "../cert/apiclient_cert.pem",
		DebugSwitch:   0,
		NotifyUrl:     "https://www.pandaball.cc",
	}
	conf := &pay.PayConfig{
		PayEngine:  "self",
		AlipayMode: "default",
		WechatMode: "default",
		Alipay:     alipay,
		Wechat:     wechat,
	}
	//request := pay.PayRequest{
	//	Amount:     1,
	//	Subject:    "测试",
	//	OutTradeNo: "2023022317223645723",
	//	TimeExpire: time.Second * 130,
	//	NotifyUrl:  "https://www.pandaball.cc",
	//	RefundUrl:  "https://www.pandaball.cc",
	//	Attach:     "member",
	//	Account:    "oPHvT5b70TJZVaaFYhNYx4Fhoivs",
	//}
	pay := pay.NewPay(conf)
	//wechatc, err := pay.Wechat().App(&request)
	//if err != nil {
	//	fmt.Println("错误")
	//}
	//wechatc, err := pay.Wechat().Mini(&request)
	//if err != nil {
	//	fmt.Println("错误")
	//}
	//fmt.Println(fmt.Sprintf("%+v", wechatc.AppletWechat), "-----------")

	//ali, err := pay.Alipay().App(&request)
	//if err != nil {
	//	fmt.Println("错误")
	//}
	//fmt.Println(fmt.Sprintf("%+v", ali.AppAlipayInfo), "-----------")

	//alipayQuery(pay)
	//wechatQuery(pay)
	wechatRefundQuery(pay)
}
func wechatRefundQuery(ipay pay.IPAY) {
	request := pay.RefundQueryRequest{
		OutRefundNo: "20231021131003433842041_359351",
	}
	ali, err := ipay.Wechat().RefundQuery(&request)
	if err != nil {
		fmt.Println("错误")
	}
	fmt.Println(fmt.Sprintf("%+v", ali), "-----------")
}
func wechatQuery(ipay pay.IPAY) {
	request := pay.QueryRequest{
		OutTradeNo: "20221018101046102774024",
	}
	ali, err := ipay.Wechat().Query(&request)
	if err != nil {
		fmt.Println("错误")
	}
	fmt.Println(fmt.Sprintf("%+v", ali), "-----------")
}
func alipayQuery(ipay pay.IPAY) {
	request := pay.QueryRequest{
		OutTradeNo: "20221114141110545428499",
	}
	ali, err := ipay.Alipay().Query(&request)
	if err != nil {
		fmt.Println("错误")
	}
	fmt.Println(fmt.Sprintf("%+v", ali), "-----------")
}

func alipayRefundQuery(ipay pay.IPAY) {
	request := pay.RefundQueryRequest{
		OutTradeNo:  "20231021101010249242197",
		OutRefundNo: "20231021101010249242197_482227",
	}
	ali, err := ipay.Alipay().RefundQuery(&request)
	if err != nil {
		fmt.Println("错误")
	}
	fmt.Println(fmt.Sprintf("%+v", ali), "-----------")
}
