package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code": "10001", // 程序中的错误码
	"message": "xxx",    // 提示信息
	"data": {},      // 数据
}
*/

type ResponseData struct {
	Code ResCode `json:"code"`
	Msg  string  `json:"message"`
	Data any     `json:"data"`
}

// httpStatusForCode 根据码所属域决定 HTTP 状态码：
// 99 域（系统错误）→ 500，其余（业务错误）→ 200。
// 这样网关/负载均衡能依据 HTTP 状态识别节点健康，实现熔断。
func httpStatusForCode(code ResCode) int {
	if code/10000 == 99 {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func writeJSON(c *gin.Context, code ResCode, msg string, data any) {
	c.JSON(httpStatusForCode(code), &ResponseData{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func Error(c *gin.Context, code ResCode) {
	writeJSON(c, code, code.Msg(), nil)
}

func Success(c *gin.Context, data any) {
	writeJSON(c, CodeSuccess, CodeSuccess.Msg(), data)
}

func ErrorWithMsg(c *gin.Context, code ResCode, msg string) {
	writeJSON(c, code, msg, nil)
}

// externalMsg 对外接口的系统错误脱敏文案：不暴露 DB/Redis/外部依赖等技术栈细节。
const externalMaskedMsg = "系统繁忙，请稍后再试"

// ErrorExternal 对外 API 专用错误响应。
// 与 Error 的区别：当 code 属于 99 域（系统错误，如 DB/缓存/外部服务故障）时，
// 一律返回脱敏文案，避免向第三方调用方泄露内部技术栈（数据库、缓存、云厂商等）。
// 业务错误（非 99 域）仍返回真实文案，因为本就该让调用方知道业务含义。
// HTTP 状态码分级（200/500）与 Error 一致，便于外部网关做熔断与重试决策。
func ErrorExternal(c *gin.Context, code ResCode) {
	msg := code.Msg()
	if code/10000 == 99 {
		msg = externalMaskedMsg
	}
	writeJSON(c, code, msg, nil)
}
