package test

import (
	"fmt"
	"github.com/info193/tools/authenticator"
	"testing"
)

func TestAuthenticator(t *testing.T) {
	// 使用方式 手机下载Authenticator app
	google := authenticator.NewAuthenticator()
	secret := google.GetSecret()
	fmt.Println("获取密钥", secret)
	params := map[string]string{
		"width":  "300",    // 生成二维码图片宽度
		"height": "300",    // 生成二维码图片高度
		"level":  "M",      // 二维码图片纠错等级 L低7%、M中15%、Q中高25%、H高30%
		"margin": "10",     // 二维码边距
		"issuer": "胡大爷", // 发行人
		"user":   "muyu",   // 账号
		"secret": secret,   // 密钥
	}
	fmt.Println("获取动态码 账号及密钥", google.SetQrCodeData(params).GetQrcodeUser())
	fmt.Println("获取二维码动态码内容", google.SetQrCodeData(params).GetQrcode())
	fmt.Println("获取二维码动态码URL", google.SetQrCodeData(params).GetQrcodeUrl())
	code := google.GetCode(secret, 30)
	fmt.Println("获取颁发验证码", code)
	fmt.Println("验证颁发验证码", google.VerifyCode(secret, code))
}
