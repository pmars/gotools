package middleware

import (
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmars/beego/logs"
	"github.com/pmars/gotools/log"
)

// 接口日志打印
func InfoLog(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path

	c.Next()

	latency := time.Since(start)
	clientIP := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	statusColor := colorForStatus(statusCode)
	methodColor := colorForMethod(method)
	var args string

	if c.Request.Header.Get("Content-Type") == "application/json" {
		b, _ := ioutil.ReadAll(c.Request.Body)
		args = string(b)
	} else {
		c.Request.ParseForm()
		args = c.Request.Form.Encode()
	}

	logs.Debug("|%s %3d %s| %6v | %s |%s %s %s %s %s",
		statusColor, statusCode, log.Reset,
		latency,
		clientIP,
		methodColor, method, log.Reset,
		path,
		args,
	)
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return log.Green
	case code >= 300 && code < 400:
		return log.White
	case code >= 400 && code < 500:
		return log.Yellow
	default:
		return log.Red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return log.Blue
	case "POST":
		return log.Cyan
	case "PUT":
		return log.Yellow
	case "DELETE":
		return log.Red
	case "PATCH":
		return log.Green
	case "HEAD":
		return log.Magenta
	case "OPTIONS":
		return log.White
	default:
		return log.Reset
	}
}
