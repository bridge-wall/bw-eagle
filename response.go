package eagle

import (
	"net/http"

	"github.com/rhonin-cd/rhonin-eagle/middleware"
	"github.com/gin-gonic/gin"
)

func JsonResponse(c *gin.Context, code int, message, result interface{}) {
	c.JSON(http.StatusOK, middleware.NewHttpResponse(code, message, result))
}

type Response struct {
	errCode int
	message interface{}
	result  interface{}
}

func NewResponse(errCode int, message, result interface{}) *Response {
	return &Response{
		errCode: errCode,
		message: message,
		result:  result,
	}
}

func NewCodeResponse(errCode int, result interface{}) *Response {
	return &Response{
		errCode: errCode,
		message: ErrMsg(errCode),
		result:  result,
	}
}

func (r *Response) GetErrCode() int {
	return r.errCode
}

func (r *Response) GetMessage() interface{} {
	return r.message
}

func (r *Response) GetResult() interface{} {
	return r.result
}

func RenderJson(c *gin.Context, resp *Response) {
	c.JSON(http.StatusOK, middleware.NewHttpResponse(resp.GetErrCode(), resp.GetMessage(), resp.GetResult()))
}

// 终止请求
func AbortRequest(c *gin.Context, resp *Response) {
	c.AbortWithStatusJSON(http.StatusOK, middleware.NewHttpResponse(resp.GetErrCode(), resp.GetMessage(), resp.GetResult()))
}
