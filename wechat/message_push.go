package wechat

import (
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	TemplateUri = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
)

type PushData interface {
	MergeMessage(msg string) map[string]*TemplateInfo
}

type MessagePush struct {
	ToUser string                   `json:"touser"`      // wechat open_id
	TmpId  string                   `json:"template_id"` //模版id
	Url    string                   `json:"url"`         // 点击推送消息的链接地址
	Data   map[string]*TemplateInfo `json:"data"`        //string就是模版的key

	accessToken *accessToken `json:"-"`
}

func GetMessagePush(wechatAppId, wechatSecret string) *MessagePush {
	msg := MessagePush{}
	msg.accessToken = GetAccessToken(wechatAppId, wechatSecret)

	return &msg
}

type pushReturn struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   int    `json:"msgid"`
}

type TemplateInfo struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// 这里需要根据自己的模板进行相应的修改，根据自己的喜好来进行调整
func (message *MessagePush) PushSimple(toUser, tmpId, msg string) error {
	data := map[string]*TemplateInfo{
		"first":    {"Demo Service ERROR", "#173177"},
		"keyword1": {"业务错误", "#173177"},
		"keyword2": {msg, "#173177"},
		"keyword3": {time.Now().Format("2006-01-02 15:04:05"), "#173177"},
		"remark":   {"服务器运行状态监控消息，请持续关注", "#173177"},
	}
	return message.Push(toUser, "", tmpId, data)
}

// 推送消息
func (message *MessagePush) Push(toUser, url, tmpId string, data map[string]*TemplateInfo) error {
	message.TmpId = tmpId
	message.Url = url
	message.Data = data
	message.ToUser = toUser

	token, err := message.accessToken.GetToken()
	if err != nil {
		fmt.Printf("get access token ERROR:%v\n", err.Error())
		return err
	}

	request := gorequest.New()
	pushReturn := pushReturn{}
	tmpUri := fmt.Sprintf(TemplateUri, token)
	_, _, errs := request.Post(tmpUri).SendStruct(message).EndStruct(&pushReturn)
	if len(errs) != 0 {
		for _, err := range errs {
			if err != nil {
				fmt.Printf("Post Template To Wechat ERROR:%v", err.Error())
			}
		}
	}

	if pushReturn.ErrCode == 40001 {
		tokenStr, err := message.accessToken.RefreshToken()

		if err != nil {
			fmt.Printf("Refresh Wechat Token ERROR:%v", err)
		} else {
			tmpUri := fmt.Sprintf(TemplateUri, tokenStr)
			_, _, errs = request.Post(tmpUri).SendStruct(message).EndStruct(&pushReturn)

			if err != nil || pushReturn.ErrCode != 0 {
				fmt.Sprintf("Post Template To Wechat ERROR:%v", err.Error())
			}
		}
	}
	if pushReturn.ErrCode != 0 {
		return fmt.Errorf(pushReturn.ErrMsg)
	}
	fmt.Printf("%#v\n", pushReturn)

	return nil
}
