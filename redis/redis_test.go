package redis

import (
	"fmt"
	"testing"
	"time"
)

var redisDao = NewRedisDao("127.0.0.1:6379", 100, false, "")
var key = "redisKey"

func TestRedisDao_Get(t *testing.T) {
	fmt.Println(redisDao.Get(key))
}

func TestRedisDao_GetSet(t *testing.T) {
	fmt.Println(redisDao.GetSet(key, 1200))
}

func TestRedisDao_Del(t *testing.T) {
	fmt.Println(redisDao.Del(key))
}

func TestRedisDao_SetNX(t *testing.T) {
	fmt.Println(redisDao.SetNX(key, 500))
}

func TestLockAndUnLock1(t *testing.T) {
	fmt.Println(redisDao.Lock(key, 1000))
	time.Sleep(time.Millisecond * 900)
	fmt.Println(redisDao.UnLock(key, 10))
}

func TestLockAndUnLock2(t *testing.T) {
	index := 0
	for i := 0; i < 2; i++ {
		go func() {
			for {
				redisDao.LockMust(key, 1000)
				fmt.Printf("%v ", index)
				index++
				redisDao.UnLock(key, 10)
			}
		}()
	}
	time.Sleep(time.Second * 10)
}
