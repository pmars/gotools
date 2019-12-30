package gotools

import (
	"crypto/md5"
	"encoding/hex"
)

/**
md5
*/
func Md5(msg string) string {
	h := md5.New()
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}
