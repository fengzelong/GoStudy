package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	CodeOK           = 0
	CodeInvalidInput = 40001
	CodeUnauthorized = 40101
	CodeForbidden    = 40301
	CodeNotFound     = 40401
	CodeConflict     = 40901
	CodeInternal     = 50001
)

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{
		Code:    CodeOK,
		Message: "ok",
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Body{
		Code:    CodeOK,
		Message: "created",
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string) {
	ErrorWithCode(c, status, codeForStatus(status), message)
}

// ErrorWithCode 返回 HTTP 状态和稳定的业务错误码，便于客户端按业务语义处理失败。
func ErrorWithCode(c *gin.Context, status int, code int, message string) {
	c.JSON(status, Body{
		Code:    code,
		Message: message,
	})
}

func codeForStatus(status int) int {
	switch status {
	case http.StatusBadRequest:
		return CodeInvalidInput
	case http.StatusUnauthorized:
		return CodeUnauthorized
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusConflict:
		return CodeConflict
	default:
		return CodeInternal
	}
}
