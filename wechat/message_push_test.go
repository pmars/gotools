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
		"wx74c099e4b07c6****",
		"4a83e84f6d61e09566****",
		"r-2ze8de57896****",
		"zheGeY6****",
		"weixi6****")

	data := map[string]*TemplateInfo{
		"first":    {"Demo Service ERROR", "#173177"},
		"keyword1": {"业务错误", "#173177"},
		"keyword2": {"这是推送测试消息", "#173177"},
		"keyword3": {time.Now().Format("2006-01-02 15:04:05"), "#173177"},
		"remark":   {"服务器运行状态监控消息，请持续关注", "#173177"},
	}

	err := message.Push(
		"oaJdSs2Y9kuRUye********",
		"http://xiaoh.me",
		"fG087aPIU794RyQy********",
		data)
	fmt.Println(err)

	err = message.PushSimple(
		"oaJdSs2Y9kuRUye********",
		"fG087aPIU794RyQyUG********-VBu4",
		"testing push message")
	fmt.Println(err)
}
