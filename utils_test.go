package gotools

import (
	"fmt"
	"testing"
)

func TestData2Str(t *testing.T) {
	data := map[string]string{
		"name": "xingming",
		"addr": "beijing",
		"work": "golang",
	}
	fmt.Println(Data2Str(data))
	fmt.Println(Data2Str("abc"))
}

func TestS2B(t *testing.T) {
	fmt.Println(S2B("abc"))
}
