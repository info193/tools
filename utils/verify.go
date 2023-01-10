package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 验证手机号码
func CheckMobile(phone string) bool {
	rule := "^1[345789]{1}\\d{9}$"
	reg := regexp.MustCompile(rule)
	return reg.MatchString(phone)
}

// 验证身份证号码
func CheckIdCard(card string) bool {
	//18位身份证 ^(\d{17})([0-9]|X)$
	// 匹配规则
	// (^\d{15}$) 15位身份证
	// (^\d{18}$) 18位身份证
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	rule := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"
	// 正则调用规则
	reg := regexp.MustCompile(rule)
	// 返回 MatchString 是否匹配
	return reg.MatchString(card)
}

// 验证邮箱
func CheckEmail(email string) bool {
	rule := "^([a-z0-9_\\.-]+)@([\\da-z\\.-]+)\\.([a-z\\.]{2,6})$"
	// 正则调用规则
	reg := regexp.MustCompile(rule)
	// 返回 MatchString 是否匹配
	return reg.MatchString(email)
}

// 验证邮箱
func CheckUrl(url string) bool {
	rule := "^(https?:\\/\\/)?([\\da-z\\.-]+)\\.([a-z\\.]{2,6})([\\/\\w \\.-]*)*\\/?$"
	// 正则调用规则
	reg := regexp.MustCompile(rule)
	// 返回 MatchString 是否匹配
	return reg.MatchString(url)
}

func SubIdCardBirthday(card string) string {
	cardStr := []rune(card)
	length := len(cardStr)
	if length == 15 {
		return fmt.Sprintf("19%v", string(cardStr[6:12]))
	}
	if length == 18 {
		return fmt.Sprintf("%v", string(cardStr[6:14]))
	}
	return ""
}

// 身份证提前性别 0未知 1男 2女
func SubIdCardGender(card string) int64 {
	cardStr := []rune(card)
	length := len(cardStr)
	if length == 15 {
		gender, _ := strconv.Atoi(string(cardStr[len(cardStr)-2 : len(cardStr)-1]))
		val := gender % 2
		if val == 0 {
			return 2
		}
		if val == 1 {
			return 1
		}
	}
	if length == 18 {
		gender, _ := strconv.Atoi(string(cardStr[len(cardStr)-2 : len(cardStr)-1]))
		val := gender % 2
		if val == 0 {
			return 2
		}
		if val == 1 {
			return 1
		}
	}
	return 0
}

// 脱敏显示
func DesensitizeString(str string, start, end int, replate string) string {
	if str == "" {
		return ""
	}
	if replate == "" {
		replate = "*"
	}
	runeStr := []rune(str)
	if len(runeStr) < 6 && start == 0 && end == 0 {
		return fmt.Sprintf("%v%v", string(runeStr[:2]), strings.Repeat(replate, len(runeStr)-2))
	}
	if len(runeStr) < 6 && start != 0 && end == 0 {
		return fmt.Sprintf("%v%v", string(runeStr[:start]), strings.Repeat(replate, len(runeStr)-start))
	}
	if len(runeStr) >= 6 && start != 0 && end == 0 {
		return fmt.Sprintf("%v%v", string(runeStr[:start]), strings.Repeat(replate, len(runeStr)-start))
	}
	if len(runeStr) >= 6 && start == 0 && end == 0 {
		return fmt.Sprintf("%v%v%v", string(runeStr[:3]), strings.Repeat(replate, len(runeStr)-3), string(runeStr[len(runeStr)-3:]))
	}
	if len(runeStr) >= 6 && start != 0 && end != 0 {
		return fmt.Sprintf("%v%v%v", string(runeStr[:start]), strings.Repeat(replate, len(runeStr)-start), string(runeStr[len(runeStr)-end:]))
	}
	return ""
}
