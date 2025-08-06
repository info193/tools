package authenticator

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"
)

type Steganography struct {
}

func NewSteganography() *Steganography {
	return &Steganography{}
}

// 将字符串转换为二进制字符串
func (l *Steganography) binaryString(s string) string {
	var binaryBuilder strings.Builder

	for i := 0; i < len(s); i++ {
		binaryBuilder.WriteString(fmt.Sprintf("%08b", s[i]))
	}
	binaryBuilder.WriteString("00000000")
	return binaryBuilder.String()
}

// 设置字节的最低有效位
func (l *Steganography) setLSB(b byte, bit byte) byte {
	// 确保 bit 是有效值
	if bit != '0' && bit != '1' {
		return b
	}
	// 清除最低位
	b = b & 0xFE
	// 设置新的最低位
	if bit == '1' {
		b |= 1
	}
	return b
}
func (l *Steganography) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 将二进制字符串转换为普通字符串
func (l *Steganography) binaryToStr(binary string) string {
	if len(binary) == 0 {
		return ""
	}

	// 确保二进制字符串长度是8的倍数
	if len(binary)%8 != 0 {
		binary = binary[:len(binary)-len(binary)%8]
	}

	text := ""
	for i := 0; i < len(binary); i += 8 {
		end := i + 8
		if end > len(binary) {
			break
		}
		byteStr := binary[i:end]

		// 将8位二进制转换为字节
		var byteVal byte
		for j := 0; j < 8; j++ {
			byteVal <<= 1
			if byteStr[j] == '1' {
				byteVal |= 1
			}
		}

		text += string(byteVal)
	}
	return text
}

func (l *Steganography) Encrypt(inputImg string, outputImg string, message string) error {
	binaryMessage := l.binaryString(message)
	messageLength := len(binaryMessage)

	file, err := os.Open(inputImg)
	if err != nil {
		return err
	}
	defer file.Close()

	// 解码PNG图像
	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	// 获取图像边界
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 检查图像是否足够大
	maxEmbed := l.min(width, height) // 对角线可嵌入的最大位数
	if maxEmbed < messageLength {
		//fmt.Printf("警告: 图像只能嵌入 %d 位，消息需要 %d 位\n", maxEmbed, messageLength)
		messageLength = maxEmbed // 截断消息以适应图像
	}

	// 创建可修改的NRGBA图像（避免预乘alpha问题）
	rgba := image.NewNRGBA(bounds)
	draw.Draw(rgba, rgba.Bounds(), img, bounds.Min, draw.Src)
	// 嵌入消息（调试模式）
	for x := 0; x < messageLength; x++ {
		y := x // 对角线坐标 (x, x)

		// 获取当前像素颜色
		c := rgba.NRGBAAt(x, y)
		r := c.R
		g := c.G
		b := c.B
		// 修改蓝色通道的LSB
		newB := l.setLSB(b, binaryMessage[x])

		// 更新像素颜色
		rgba.SetNRGBA(x, y, color.NRGBA{
			R: r,
			G: g,
			B: newB,
			A: c.A,
		})
	}

	// 保存修改后的图像
	outFile, err := os.Create(outputImg)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 使用默认PNG编码器（避免压缩修改）
	err = png.Encode(outFile, rgba)
	if err != nil {
		return err
	}
	return nil
}

func (l *Steganography) Decrypt(src string) (string, error) {
	// 打开PNG文件
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// 解码PNG图像
	img, err := png.Decode(f)
	if err != nil {
		return "", err
	}

	// 获取图像边界
	bounds := img.Bounds()
	// 创建可修改的NRGBA图像（避免预乘alpha问题）
	rgba := image.NewNRGBA(bounds)
	draw.Draw(rgba, rgba.Bounds(), img, bounds.Min, draw.Src)

	binMessage := ""
	x := 0
	for {
		y := x // 使用对角线上的像素 (x, x)
		// 获取当前像素颜色
		c := rgba.NRGBAAt(x, y)

		// 将蓝色通道值转换为8位二进制字符串
		blueBin := fmt.Sprintf("%08b", c.B)
		// 取最低位（LSB）
		binMessage += blueBin[7:]

		// 检查是否遇到终止符"00000000"
		if len(binMessage) >= 8 {
			termination := binMessage[len(binMessage)-8:]
			if termination == "00000000" {
				binMessage = binMessage[:len(binMessage)-8]
				break
			}
		}
		x++
	}
	//fmt.Println(binMessage)
	// 转换二进制字符串为普通字符串
	return l.binaryToStr(binMessage), nil
}
