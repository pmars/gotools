package wechat

import (
	"fmt"
	"testing"
	"time"
)

/*
	jsonStr := `{
		"appId":"wxa4af6593684a4f0b",
		"secret":"2b2cc3577ff534ae2bc7c54****",
		"level":3,
		"tmpId":"IDkJ_4ulcyWgz0flVvwXzlqsvkenbnMV4kU****",
		"wxIds":["oDncy0udgxSpnL1Dr7TJWLyRI6_g"]
	}`
	log.InitWechatPush(jsonStr, initPushData)
*/
func TestMessagePush_Push(t *testing.T) {
	message := GetMessagePush(
		"wx74c09********",
		"4a83e8********",
		"********",
		"********",
		"********")

	data := map[string]*TemplateInfo{
		"first":    {"Demo Service ERROR", "#173177"},
		"keyword1": {"业务错误", "#173177"},
		"keyword2": {"这是推送测试消息", "#173177"},
		"keyword3": {time.Now().Format("2006-01-02 15:04:05"), "#173177"},
		"remark":   {"服务器运行状态监控消息，请持续关注", "#173177"},
	}

	err := message.Push(
		"oaJdSs2Y********",
		"http://xiaoh.me",
		"fG087aPIU7********",
		data)
	fmt.Println(err)

	err = message.PushSimple(
		"oaJdSs2********",
		"fG087aPIU7********",
		"testing push message")
	fmt.Println(err)
}
