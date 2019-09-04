package log

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"fmt"

	"github.com/pmars/beego/logs"
	"github.com/pmars/gotools"
	"github.com/pmars/gotools/wechat"
)

var (
	AdapterWechatPush = "wechat_push"
	InitPushData      func(string) map[string]*wechat.TemplateInfo
)

type WechatPush struct {
	AppId     string   `json:"appId"`
	Secret    string   `json:"secret"`
	RedisConn string   `json:"redisConn"`
	RedisAuth string   `json:"redisAuth"`
	RedisKey  string   `json:"redisKey"`
	Level     int      `json:"level"`
	HostName  string   `json:"host"`
	OuterIp   string   `json:"ip"`
	TmpId     string   `json:"tmpId"` //微信消息模版ID
	WechatIds []string `json:"wxIds"` //关注微信公众号，接受错误推送者的openid集合

	message *wechat.MessagePush
}

/*
	初始化微信
*/
func InitWechatPush(jsonConf string, initPushData func(string) map[string]*wechat.TemplateInfo) error {
	InitPushData = initPushData
	logs.Register(AdapterWechatPush, NewWechatPush)
	return logs.SetLogger(AdapterWechatPush, jsonConf)
}

// 返回一个 logs.Logger 的接口
func NewWechatPush() logs.Logger {
	wechatPush := &WechatPush{}
	wechatPush.HostName, _ = os.Hostname()
	wechatPush.OuterIp = gotools.GetOutboundIP().String()

	return wechatPush
}

// Logger的初始化函数，conf是在setLogger的时候传入的config字符串
func (wechatPush *WechatPush) Init(conf string) error {
	err := json.Unmarshal([]byte(conf), wechatPush)
	if err != nil {
		fmt.Printf("Init Wechat Push:%v System ERROR:%v\n", gotools.Data2Str(wechatPush), err.Error())
		return err
	}

	if len(wechatPush.TmpId) == 0 || len(wechatPush.WechatIds) == 0 ||
		len(wechatPush.AppId) == 0 || len(wechatPush.Secret) == 0 {
		return errors.New("wechat push args error")
	}

	wechatPush.message = wechat.GetMessagePush(wechatPush.AppId, wechatPush.Secret, wechatPush.RedisConn,
		wechatPush.RedisAuth, wechatPush.RedisKey)

	return nil
}

// WriteMsg will write the msg and level
func (wechatPush *WechatPush) WriteMsg(when time.Time, msg string, level int) error {
	if level > wechatPush.Level {
		return nil
	}

	data := InitPushData(msg)

	for _, id := range wechatPush.WechatIds {
		err := wechatPush.message.Push(id, "", wechatPush.TmpId, data)
		fmt.Printf("push data to user:%v, error:%v\n", id, err)
	}
	return nil
}

// Destroy is a empty method
func (wechatPush *WechatPush) Destroy() {

}

// Flush is a empty method
func (wechatPush *WechatPush) Flush() {

}
