package wechat

import (
	"fmt"
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
	pRedis "github.com/pmars/gotools/redis"

	"sync"

	"github.com/parnurzeal/gorequest"
	"github.com/pmars/gotools"
)

const (
	WxGetTokenUrl = "https://api.weixin.qq.com/cgi-bin/token?"
)

type accessToken struct {
	wechatAppId  string
	wechatSecret string
	redisObj     *pRedis.RedisDao

	_accessToken  string
	redisTokenKey string

	expiresIn time.Time
	mutex     sync.RWMutex
}

func GetAccessToken(appId, secret, redisConn, redisAuth, redisKey string) *accessToken {
	accessP := &accessToken{
		wechatAppId:  appId,
		wechatSecret: secret,
	}
	if redisKey != "" {
		accessP.redisTokenKey = redisKey
		accessP.redisObj = pRedis.NewRedisDao(redisConn, 10, true, redisAuth)
	}
	return accessP
}

type wechatToken struct {
	AccessToken string `json:"access_token" redis:"access_token"`
	ExpiresIn   int    `json:"expires_in" redis:"expires_in"`
	ExpiresTime int    `json:"expires_time" redis:"expires_time"`
}

// 网络请求微信服务，获取AccessToken
func (token *accessToken) RefreshToken() (string, error) {
	token.mutex.Lock()
	defer token.mutex.Unlock()

	v := url.Values{}
	v.Set("appid", token.wechatAppId)   //公众号appid
	v.Set("secret", token.wechatSecret) //公众号secret
	v.Set("grant_type", "client_credential")

	wt := wechatToken{}
	tokenUrl := WxGetTokenUrl + v.Encode()

	if _, body, errs := gorequest.New().Get(tokenUrl).EndStruct(&wt); errs != nil || len(wt.AccessToken) == 0 {
		return "", fmt.Errorf("get access token error:%v, tokenUrl:%v result:%v", errs, tokenUrl, string(body))
	}

	token.setToken(wt.AccessToken, wt.ExpiresIn)
	token.expiresIn = time.Now().Add(time.Second * time.Duration(wt.ExpiresIn))

	fmt.Printf("token:%v, expires:%v\n", wt.AccessToken, token.expiresIn.Format(gotools.TimeFormat))
	return wt.AccessToken, nil
}

// 使用缓存，获取AccessToken
func (token *accessToken) GetToken() (string, error) {
	token.mutex.RLock()

	if token.expiresIn.After(time.Now()) {
		token.mutex.RUnlock()
		return token.getToken()
	}
	token.mutex.RUnlock()

	return token.RefreshToken()
}

func (token *accessToken) getToken() (string, error) {
	// 如果没有用到Redis存储，则判断是否有_accessToken字段
	fmt.Println("token.redisTokenKey", token.redisTokenKey)
	if token.redisTokenKey == "" {
		if token._accessToken != "" {
			// 如果有，则直接返回
			return token._accessToken, nil
		} else {
			// 如果没有，则刷新该字段
			return token.RefreshToken()
		}
	}

	aToken := wechatToken{}
	conn := token.redisObj.RedisPool().Get()
	defer conn.Close()

	v, err := redis.Values(conn.Do("HGETALL", token.redisTokenKey))
	if err == nil {
		err = redis.ScanStruct(v, &aToken)
	} else {
		fmt.Println("redis HGETALL ERROR:", err.Error())
	}

	if pRedis.RedisOK(err) == nil {
		return aToken.AccessToken, nil
	}

	return token.RefreshToken()
}

func (token *accessToken) setToken(tk string, expiresIn int) error {
	if token.redisTokenKey == "" {
		token._accessToken = tk
	} else {
		conn := token.redisObj.RedisPool().Get()
		defer conn.Close()

		aToken := wechatToken{
			AccessToken: tk,
			ExpiresIn:   expiresIn,
			ExpiresTime: expiresIn + int(time.Now().Unix()),
		}
		_, err := redis.String(conn.Do("HMSET", redis.Args{}.Add(token.redisTokenKey).AddFlat(aToken)...))
		return pRedis.RedisOK(err)
	}
	return nil
}

func (token *accessToken) _getToken() (string, error) {
	if token.redisTokenKey == "" && token._accessToken != "" {
		// 如果没有用到Redis存储，则判断是否有_accessToken字段，如果有，则直接返回
		return token._accessToken, nil
	} else if token.redisTokenKey == "" && token._accessToken == "" {
		// 如果没有用到Redis存储，则判断是否有_accessToken字段，如果没有，则刷新该字段
		return token.RefreshToken()
	} else if tk, err := token.redisObj.Get(token.redisTokenKey); err != nil {
		// 使用Redis的方案，获取的时候出问题了，则刷新一下
		return token.RefreshToken()
	} else {
		// 成功从Redis里面获取数据，则直接返回
		return tk, nil
	}
}

func (token *accessToken) _setToken(tk string, expiresIn int) error {
	if token.redisTokenKey == "" {
		token._accessToken = tk
	} else {
		conn := token.redisObj.RedisPool().Get()
		defer conn.Close()

		_, err := conn.Do("SET", token.redisTokenKey, tk)
		redis.Int(conn.Do("EXPIREAT", token.redisTokenKey, expiresIn))

		return pRedis.RedisOK(err)
	}
	return nil
}
