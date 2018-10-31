package redis

import "time"
import (
	"fmt"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

const REDISOK = "redigo: nil returned" //redis正常返回
func RedisOK(err error) (err2 error) {
	if err == nil || err.Error() == REDISOK {
		return nil
	} else {
		return err
	}
}

type RedisDao struct {
	redisPool *redis.Pool
}

func NewRedisDao(server string, maxConn int, isAuth bool, auth string) *RedisDao {
	r := &RedisDao{}
	if isAuth {
		r.redisPool = GetRedisPoolAuth(server, auth, maxConn)
	} else {
		r.redisPool = GetRedisPool(server, maxConn)
	}
	return r
}

func (dao *RedisDao) RedisPool() *redis.Pool {
	return dao.redisPool
}

func GetRedisPool(server string, maxConn int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     200,
		MaxActive:   maxConn,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func GetRedisPoolAuth(server, auth string, maxConn int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     200,
		MaxActive:   maxConn,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", auth); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// Get
func (dao *RedisDao) Get(key string) (s string, err error) {
	conn := dao.redisPool.Get()
	defer conn.Close()
	s, err = redis.String(conn.Do("GET", key))
	err = RedisOK(err)
	return
}

// GetSet
func (dao *RedisDao) GetSet(key string, value interface{}) (s string, err error) {
	conn := dao.redisPool.Get()
	defer conn.Close()
	s, err = redis.String(conn.Do("GETSET", key, value))
	err = RedisOK(err)
	return s, nil
}

// Del 可以删除多个key 返回删除key的num和错误
func (dao *RedisDao) Del(key ...interface{}) (num int, err error) {
	conn := dao.redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("DEL", key...))
	err = RedisOK(err)
	return
}

// 不存在则设置，存在则不设置
func (dao *RedisDao) SetNX(key string, value interface{}) (success bool, err error) {
	conn := dao.redisPool.Get()
	defer conn.Close()
	num, err := redis.Int(conn.Do("SETNX", key, value))
	err = RedisOK(err)
	return num > 0, err
}

// 所有的过期时间都是毫秒（ms）单位
func getExpiredTime(timeout int) (int64, int64) {
	now := time.Now()
	return now.UnixNano(), now.Add(time.Millisecond * time.Duration(timeout)).UnixNano()
}

func (dao *RedisDao) Lock(key string, timeout int) (locked bool, expiredTime int64, err error) {
	// 根据过期时间毫秒数获取当前时间和过期时间
	now, expiredTime := getExpiredTime(timeout)

	// SetNX 设置过期时间，如果成功，加锁成功，如果失败，证明锁被占据
	if success, err := dao.SetNX(key, expiredTime); err != nil {
		return false, 0, fmt.Errorf("SetNX key:%v, value:%v, ERROR:%v", key, expiredTime, err)
	} else if success {
		return true, expiredTime, nil
	}

	// 如果锁被占据，获取锁的内容，判断是否过期等
	if value, err := dao.Get(key); err != nil {
		return false, 0, fmt.Errorf("get key:%v, ERROR:%v", key, err)
	} else if value == "" {
		// 已经被删除的情况下，返回的是空
		if value, err := dao.GetSet(key, expiredTime); err != nil {
			return false, 0, fmt.Errorf("GetSet key:%v, expireTime:%v, ERROR:%v", key, expiredTime, err)
		} else if value == "" {
			return true, expiredTime, nil
		}
	} else {
		// 如果里面有内容，查看是否过期即可
		if passTime, err := strconv.ParseInt(value, 10, 64); err != nil {
			// 解析时出现问题，这个时候问题较大，只能等待过期时间了（否则在这边删除也可以，不过正常情况系啊，不会走这边）
			return false, 0, fmt.Errorf("ParseInt passTime:%v ERROR:%v", value, err)
		} else if now > passTime {
			// 已经过期了，这个时候直接获取锁即可
			if valueNow, err := dao.GetSet(key, expiredTime); err != nil {
				return false, 0, fmt.Errorf("GetSet key:%v, expireTime:%v, ERROR:%v", key, expiredTime, err)
			} else if valueNow == value {
				return true, expiredTime, nil
			}
		}
	}
	return false, 0, nil
}

func (dao *RedisDao) LockRetry(key string, timeout, retryTimes int) (locked bool, expiredTime int64, err error) {
	for i := 0; i < retryTimes; i++ {
		if locked, expiredTime, err = dao.Lock(key, timeout); err != nil || locked {
			return
		}
		time.Sleep(time.Millisecond * time.Duration(i+1))
	}
	return false, 0, nil
}

func (dao *RedisDao) LockMust(key string, timeout int) (locked bool, expiredTime int64, err error) {
	for i := 0; ; i++ {
		if locked, expiredTime, err = dao.Lock(key, timeout); err != nil || locked {
			return
		}
		time.Sleep(time.Millisecond * time.Duration(i+1))
	}
}

func (dao *RedisDao) UnLock(key string, safeDelTime int64) (bool, error) {
	if value, err := dao.Get(key); err != nil {
		// 获取KEY的时候报错，证明可能已经过期，或者别别人删除了
		return false, fmt.Errorf("get key:%v, ERROR:%v", key, err)
	} else if expireTime, err := strconv.ParseInt(value, 10, 64); err != nil {
		// 过期时间解析错误
		return false, fmt.Errorf("ParseInt key:%v, ERROR:%v", value, err)
	} else if time.Now().UnixNano()+safeDelTime*1000000 > expireTime {
		// 就要到过期时间了，这个时候直接退出，等待过期时间即可
		return false, fmt.Errorf(" Key:%v nearly to be expired.", key)
	} else if count, err := dao.Del(key); err != nil {
		return false, fmt.Errorf(" Del key:%v, ERROR:%v", key, err)
	} else if count == 0 {
		return false, fmt.Errorf(" Del key:%v, Count:%v", key, count)
	}
	return true, nil
}
