package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5(s string) string {
	MD5 := md5.New()
	_, _ = io.WriteString(MD5, s)
	return hex.EncodeToString(MD5.Sum(nil))
}
