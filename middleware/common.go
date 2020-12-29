package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpResponse struct {
	ErrCode int         `json:"error_code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func NewHttpResponse(code int, message, result interface{}) *HttpResponse {
	resp := &HttpResponse{}
	resp.ErrCode = code

	if message == nil {
		resp.Message = ""
	} else {
		resp.Message = fmt.Sprint(message)
	}

	if result == nil {
		resp.Result = struct{}{}
	} else {
		resp.Result = result
	}

	return resp
}

// 发生错误终止请求
func abortRequest(c *gin.Context, code int, message string, result interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, NewHttpResponse(code, message, result))
}
