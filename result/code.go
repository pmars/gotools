package result

import "github.com/gin-gonic/gin"

const (
	CodeKey = "returnCode"
)

var (
	Success = ReturnCode{0, "Success", "Success"}

	ErrArgs            = ReturnCode{6002, "参数错误", "Args Error"}
	ErrNoData          = ReturnCode{6012, "没有查询记录", "Args Error"}
	ErrPass            = ReturnCode{6004, "密码错啦", "Password Error"}
	ErrUserDisabled    = ReturnCode{6011, "用户已被禁用", "User Disabled"}
	ErrUserRoleLimited = ReturnCode{6101, "用户角色错误", "User Role Error"}
	ErrSQLErr          = ReturnCode{6013, "数据获取失败", "data Error"}

	ErrServer = ReturnCode{50000, "网络开小差了，一会儿再来试试吧~", "Server Error"}
)

// 此方法主要在Controller中使用，设置错误码专用
func SetCode(c *gin.Context, code ReturnCode) {
	c.Set(CodeKey, code)
}

// 在返回结果之前调用，获取错误码，如果没有设置，默认成功
func GetCode(c *gin.Context) ReturnCode {
	code, exist := c.Get(CodeKey)
	if !exist {
		return Success
	}
	return code.(ReturnCode)
}

type ReturnCode struct {
	Code  int
	CnMsg string
	EnMsg string
}
