package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


// Http Header
const (
	HttpHeaderTraceId      = "Trace-ID"
)

// 请求上下文
const (
	CtxRequestStart = "bw_request_start" // 请求开始时间
	CtxRequestSize  = "bw_request_size"  // 请求大小
	CtxRequestBody  = "bw_request_body"  // 请求Body
)

// 日志记录请求响应数据大小上限，太长影响性能
const (
	LogRequestMaxSize  = 1024
	LogResponseMaxSize = 1024
)

// 初始化注入上下文参数
// gin.Context.Keys类型是map，不保证线程安全
// 在多个goroutine共享，存在冲突的风险
// https://github.com/gin-gonic/gin/issues/700
// 如果context需要在多个goroutine共享，使用context.Copy()
// 【备注】新的版本已通过读写锁 RWMutex 解决了这个问题
func HandlerInitContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		setCtxRequestStart(c)
		setCtxRequestBody(c)
		setCtxRequestSize(c)

		c.Next()
	}
}

func GetCtxParam(c *gin.Context, key string) interface{} {
	value, exists := c.Get(key)
	if !exists {
		return nil
	}

	return value
}

func GetCtxParamInt(c *gin.Context, key string) (i int) {
	value := GetCtxParam(c, key)
	if value != nil {
		if v, ok := value.(int); ok {
			i = v
		}
	}
	return
}

func GetCtxParamString(c *gin.Context, key string) (s string) {
	value := GetCtxParam(c, key)
	if value != nil {
		if v, ok := value.(string); ok {
			s = v
		}
	}
	return
}

func GetCtxParamTime(c *gin.Context, key string) (t time.Time) {
	value := GetCtxParam(c, key)
	if value != nil {
		if v, ok := value.(time.Time); ok {
			t = v
		}
	}
	return
}

// 请求开始时间
func setCtxRequestStart(c *gin.Context) {
	c.Set(CtxRequestStart, time.Now())
}

// 请求数据
func setCtxRequestBody(c *gin.Context) {
	reqBody, _ := c.GetRawData()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
	body := string(reqBody)
	c.Set(CtxRequestBody, body)
}

// 请求数据大小
func setCtxRequestSize(c *gin.Context) {
	s := computeApproximateRequestSize(c.Request)
	c.Set(CtxRequestSize, s)
}

// From https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}

	return s
}

