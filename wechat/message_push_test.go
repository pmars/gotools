package wechat

import (
	"fmt"
	"testing"
	"time"
)

func TestMessagePush_Push(t *testing.T) {
	message := GetMessagePush("wx3f3a43ee700b****", "aa8b08b52ed52fe89e53a181****")

	data := map[string]*TemplateInfo{
		"first":    {"Demo Service ERROR", "#173177"},
		"keyword1": {"业务错误", "#173177"},
		"keyword2": {"这是推送测试消息", "#173177"},
		"keyword3": {time.Now().Format("2006-01-02 15:04:05"), "#173177"},
		"remark":   {"服务器运行状态监控消息，请持续关注", "#173177"},
	}

	err := message.Push(
		"oKdl5vz4q_YL9P-VTvUhTXsBrEq8",
		"http://xiaoh.me",
		"GZuiwMgwXdDyJduA33op5rhf2svf0uOH5N4dxx8Il0Q",
		data)
	fmt.Println(err)

	err = message.PushSimple(
		"oKdl5vz4q_YL9P-VTvUhTXsBrEq8",
		"GZuiwMgwXdDyJduA33op5rhf2svf0uOH5N4dxx8Il0Q",
		"testing push message")
	fmt.Println(err)
}
