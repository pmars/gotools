package wechat

import (
	"fmt"
	"net/url"
	"time"

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

	accessToken string
	expiresIn   time.Time
	mutex       sync.RWMutex
}

func GetAccessToken(appId, secret string) *accessToken {
	return &accessToken{
		wechatAppId:  appId,
		wechatSecret: secret,
	}
}

type wechatToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
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

	if _, _, errs := gorequest.New().Get(tokenUrl).EndStruct(&wt); errs != nil || len(wt.AccessToken) == 0 {
		return "", fmt.Errorf("get access token error:%v, accessToken:%v", errs, wt.AccessToken)
	}

	token.accessToken = wt.AccessToken
	token.expiresIn = time.Now().Add(time.Second * time.Duration(wt.ExpiresIn))

	fmt.Printf("token:%v, expires:%v\n", token.accessToken, token.expiresIn.Format(gotools.TimeFormat))
	return token.accessToken, nil
}

// 使用缓存，获取AccessToken
func (token *accessToken) GetToken() (string, error) {
	token.mutex.RLock()

	if token.expiresIn.After(time.Now()) {
		token.mutex.RUnlock()
		return token.accessToken, nil
	}
	token.mutex.RUnlock()

	return token.RefreshToken()
}
