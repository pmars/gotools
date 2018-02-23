package gotools

import (
	"unsafe"

	"fmt"

	"github.com/gin-gonic/gin/json"
)

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
