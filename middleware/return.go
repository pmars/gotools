package middleware

import (
	"github.com/pmars/gotools/result"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmars/beego/logs"
	"github.com/pmars/gotools"
)

// 通用响应处理中间件
func SetReturn(c *gin.Context) {
	// 调用请求
	c.Next()

	// 获取Controller层设置的Return信息
	rst := result.GetResultMap(c)
	rst.SetReturnMsg(c)

	logs.Debug("Middle Return: %v", gotools.Data2Str(rst))
	c.JSON(http.StatusOK, rst)
}
