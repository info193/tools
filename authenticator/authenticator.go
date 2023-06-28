package authenticator

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"net/url"
	"time"
)

type Authenticator struct {
	width  string
	height string
	margin string
	issuer string
	level  string
	user   string
	secret string
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{
		width:  "200",
		height: "200",
		level:  "M",
		issuer: "",
		margin: "10", // 二维码边距
		user:   "",
		secret: "",
	}
}
func (l *Authenticator) GetSecret() string {
	randomStr := l.randStr(16)
	return randomStr
}
func (l *Authenticator) randStr(strSize int) string {
	dictionary := "KDW7J2PTQAZUCR4YBFMVOX3IHLGS6NE5"
	var bytes = make([]byte, strSize)
	_, _ = rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
func (l *Authenticator) SetQrCodeData(params map[string]string) *Authenticator {
	if value, ok := params["width"]; ok {
		l.width = value
	}
	if value, ok := params["height"]; ok {
		l.height = value
	}
	if value, ok := params["issuer"]; ok {
		l.issuer = value
	}
	if value, ok := params["level"]; ok {
		l.level = value
	}
	if value, ok := params["user"]; ok {
		l.user = value
	}
	if value, ok := params["secret"]; ok {
		l.secret = value
	}
	// 二维码边距
	if value, ok := params["margin"]; ok {
		l.margin = value
	}
	return l
}
func (l *Authenticator) SetSecret(secret string) *Authenticator {
	l.secret = secret
	return l
}

// 获取动态码二维码内容
func (l *Authenticator) GetQrcode() string {
	if l.user == "" || l.secret == "" {
		return ""
	}
	issuer := ""
	if l.issuer != "" {
		issuer = fmt.Sprintf("&issuer=%s", l.issuer)
	}
	return fmt.Sprintf("otpauth://totp/%s?secret=%s%s", l.user, l.secret, issuer)
}

// 获取动态码账号
func (l *Authenticator) GetQrcodeUser() map[string]string {
	if l.user == "" || l.secret == "" {
		return nil
	}
	auth := make(map[string]string, 0)
	auth["username"] = fmt.Sprintf("%s:%s", l.user, l.issuer)
	auth["secret"] = l.secret
	return auth
}

// 获取动态码二维码url
func (l *Authenticator) GetQrcodeUrl() string {
	urlencode := url.QueryEscape(l.GetQrcode())
	if urlencode == "" {
		return ""
	}
	// 接口文档地址 https://goqr.me/api/doc/create-qr-code/
	// data= 二维码内容
	// size= 生成二维码图片 宽高
	// ecc= 代表纠错水平 L（低，大约7%的破坏数据可以纠正）M（中间，大约15%的破坏数据可以纠正）Q（质量，大约25%的破坏数据可以纠正）H（高，大约30%的破坏数据可以纠正）
	return fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?data=%s&size=%sx%s&ecc=%s&margin=%s", urlencode, l.width, l.height, l.level, l.margin)

	////    &cht=qr 这是说图表类型为qr也就是二维码。
	////    &chs=200×200 这是说生成图片尺寸为200×200，是宽x高。这并不是生成图片的真实尺寸，应该是最大尺寸吧。
	////    &choe=UTF-8 这是说内容的编码格式为UTF-8，此值默认为UTF-8.其他的编码格式请参考Google API文档。
	////    &chld=L|4 L代表默认纠错水平； 4代表二维码边界空白大小，可自行调节。具体参数请参考Google API文档。
	////    &chl=XXXX 这是QR内容，也就是解码后看到的信息。包含中文时请使用UTF-8编码汉字，否则将出现问题。
	//return fmt.Sprintf("https://chart.googleapis.com/chart?chs=%sx%s&chld=%s|0&cht=qr&chl=%s", width, height, level, urlencode)
}

// 为了考虑时间误差，判断前当前时间及前后30秒时间
func (l *Authenticator) VerifyCode(secret string, code int32) bool {
	// 当前google值
	if l.GetCode(secret, 0) == code {
		return true
	}

	// 前30秒google值
	if l.GetCode(secret, -30) == code {
		return true
	}

	// 后30秒google值
	if l.GetCode(secret, 30) == code {
		return true
	}

	return false
}

// 获取Google Code
func (l *Authenticator) GetCode(secret string, offset int64) int32 {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix() + offset
	return int32(l.oneTimePassword(key, l.toBytes(epochSeconds/30)))
}

func (l *Authenticator) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (l *Authenticator) toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func (l *Authenticator) oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := l.toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}
