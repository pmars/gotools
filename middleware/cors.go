package middleware

import "github.com/gin-gonic/gin"

// 允许服务器跨域访问
func Cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "OPTIONS, HEAD, GET, POST, DELETE, PUT")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Range, Content-Disposition")
}
