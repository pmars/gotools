package gotools

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/satori/go.uuid"
)

func MixFile(inPath, outPath string) error {
	// 读取文件
	reader, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	// 生成混淆code
	code := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	bytes = append(bytes, []byte(code)...)

	// 输出到文件
	_ = os.Remove(outPath)
	err = ioutil.WriteFile(outPath, bytes, 0777) // 写入文件内容，0777是创建的文件权限
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
