package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestCheckEmail(t *testing.T) {
	fmt.Println(utils.CheckEmail("416@163.com"))
}

func TestCheckMobile(t *testing.T) {
	fmt.Println(utils.CheckMobile("18072177633"))
}

func TestCheckIdCard(t *testing.T) {
	fmt.Println(utils.CheckIdCard("362310199011236"))
}

func TestCheckUrl(t *testing.T) {
	fmt.Println(utils.CheckUrl("http://www.baidu.com"))
}

func TestDesensitizeString(t *testing.T) {
	fmt.Println(utils.DesensitizeString("http://w哈w1w.baidu.com", 15, 3, "*"))
}

func TestSubIdCardBirthday(t *testing.T) {

	fmt.Println(utils.SubIdCardBirthday("362330861223668"))
	fmt.Println(utils.SubIdCardBirthday("421087198606017328"))

}

func TestSubIdCardGender(t *testing.T) {
	sw := utils.SubIdCardGender("362330861223668")
	if sw == 1 {
		fmt.Println("15位身份证提取出【男】", sw)
	}
	if sw == 2 {
		fmt.Println("15位身份证提取出【女】", sw)
	}

	sb := utils.SubIdCardGender("421087198606017318")
	if sb == 1 {
		fmt.Println("18位身份证提取出【男】", sb)
	}
	if sb == 2 {
		fmt.Println("18位身份证提取出【女】", sb)
	}
}
