package log

import (
	"testing"

	"time"

	"github.com/astaxie/beego/logs"
	"github.com/pmars/gotools/wechat"
)

func initPushData(msg string) map[string]*wechat.TemplateInfo {
	return map[string]*wechat.TemplateInfo{
		"first":    {"Demo Service ERROR", "#173177"},
		"keyword1": {"业务错误", "#173177"},
		"keyword2": {msg, "#173177"},
		"keyword3": {time.Now().Format("2006-01-02 15:04:05"), "#173177"},
		"remark":   {"服务器运行状态监控消息，请持续关注", "#173177"},
	}
}

func TestBeegoLogPush(t *testing.T) {
	jsonStr := `{
		"appId":"wx3f3a43ee70****",
		"secret":"aa8b08b52ed52fe89e53a1****",
		"level":3,
		"tmpId":"GZuiwMgwXdDyJduA33op5rhf2svf0uO****",
		"wxIds":["oKdl5vz4q_YL9P-VTvUhT****"]
	}`
	InitWechatPush(jsonStr, initPushData)

	logs.Error("hello error")
	logs.Debug("hello debug")
}
