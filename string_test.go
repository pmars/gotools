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

func TestRandStringRunes(t *testing.T) {
	fmt.Println(RandStringRunes(10))
	fmt.Println(RandStringRunes(8))
	fmt.Println(RandStringRunes(4))
	fmt.Println(RandStringRunes(35))
}
