package gotools

import (
	"fmt"
	"testing"
)

var inPath = "/Users/xiaoh/Downloads/1.png"
var outPath = "/Users/xiaoh/Downloads/out.png"

func Test_MixFile(t *testing.T) {
	fmt.Println(MixFile(inPath, outPath))
}
