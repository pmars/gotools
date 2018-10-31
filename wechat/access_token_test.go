package wechat

import (
	"fmt"
	"sync"
	"testing"
)

func TestAccessToken(t *testing.T) {
	wait := sync.WaitGroup{}
	access := GetAccessToken("wx3f3a43ee70***", "1e97654e0560535bbb0***")
	for i := 0; i < 100000; i++ {
		wait.Add(1)
		go func() {
			if _, err := access.GetToken(); err != nil {
				fmt.Println(err)
			}
			wait.Done()
		}()
	}

	for i := 0; i < 100; i++ {
		wait.Add(1)
		go func() {
			if _, err := access.RefreshToken(); err != nil {
				fmt.Println(err)
			}
			wait.Done()
		}()
	}

	wait.Wait()
}

func TestAccessTokenOnce(t *testing.T) {
	access := GetAccessToken("wx28ceb68f61****", "d46d03e757125776beb****")
	token, err := access.GetToken()
	fmt.Println(token)
	fmt.Println(err)
}
