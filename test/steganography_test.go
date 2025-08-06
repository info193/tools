package test

import (
	"fmt"
	"github.com/info193/tools/authenticator"
	"testing"
)

func TestSteganography(t *testing.T) {
	//encrypt() // 注意图片周边可能有大块透明度区域，应避免该区域，否则前端将无法读取。建议全图区域只有小块或几乎无透明区
	decrypt()
}
func encrypt() {
	steganography := authenticator.NewSteganography()
	err := steganography.Encrypt("800.png", "encrypt.png", "78sdfwe7r18H")
	if err != nil {
		fmt.Println("图片加密失败")
		return
	}
	fmt.Println("隐写成功")
}
func decrypt() {
	steganography := authenticator.NewSteganography()
	str, err := steganography.Decrypt("encrypt.png")
	if err != nil {
		fmt.Println("图片加密失败")
		return
	}
	fmt.Println("读取隐写数据", str)
}
