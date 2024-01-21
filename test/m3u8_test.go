package test

import (
	m3u82 "github.com/info193/tools/m3u8"
	"testing"
)

func TestM3u8(t *testing.T) {
	// 初始化设置输出路径1
	m3u8 := m3u82.NewM3U8("D:\\project\\tools/m3u8")
	// 下载m3u8文件
	urls := []string{
		"https://billiards-test.oss-cn-hangzhou.aliyuncs.com/upload/temp/demo.m3u8",
		"https://billiards-test.oss-cn-hangzhou.aliyuncs.com/upload/temp/demo1.m3u8",
		"https://billiards-test.oss-cn-hangzhou.aliyuncs.com/upload/temp/demo2.m3u8",
		"https://billiards-test.oss-cn-hangzhou.aliyuncs.com/upload/temp/demo3.m3u8",
	}
	m3u8.Download(urls)
}

