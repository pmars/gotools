package wechat

import (
	"fmt"
	"testing"
	"time"
)

func TestMessagePush_Push(t *testing.T) {
	message := GetMessagePush("wx3f3a43ee7***", "1e97654e0560535bbb0a4***")

	data := map[string]*TemplateInfo{
		"first":    {"Demo Service ERROR", "#173177"},
		"keyword1": {"业务错误", "#173177"},
		"keyword2": {"这是推送测试消息", "#173177"},
		"keyword3": {time.Now().Format("2006-01-02 15:04:05"), "#173177"},
		"remark":   {"服务器运行状态监控消息，请持续关注", "#173177"},
	}

	err := message.Push(
		"oKdl5vz4q_YL9P-VTvUhTX***",
		"http://xiaoh.me",
		"GZuiwMgwXdDyJduA33op5rhf2svf0uOH***",
		data)
	fmt.Println(err)

	err = message.PushSimple(
		"oKdl5vz4q_YL9P-VTvUhTX***",
		"GZuiwMgwXdDyJduA33op5rhf2svf0uOH***",
		"testing push message")
	fmt.Println(err)
}
