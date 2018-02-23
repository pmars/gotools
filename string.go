package gotools

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"unsafe"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Data2Str(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return string(bytes)
}

func S2B(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
