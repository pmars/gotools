package gotools

import (
	"testing"
	"fmt"
)

func TestGetNowStr(t *testing.T) {
	fmt.Println(GetNowStr())
}

func TestGetNowNoSepStr(t *testing.T) {
	fmt.Println(GetNowNoSepStr())
}

func TestGetDateStr(t *testing.T) {
	fmt.Println(GetDateStr())
}

func TestGetDateNoSepStr(t *testing.T) {
	fmt.Println(GetDateNoSepStr())
}