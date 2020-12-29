package middleware

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ngx访问日志中间件
func HandlerAccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 响应数据接口
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 打印请求日志
		ngxLogInfo(c, "")

		c.Next()

		// respLog := ""
		// body := blw.body.String()
		// if len(body) < LogResponseMaxSize {
		// 	respLog = body
		// } else {
		// 	respLog = string([]rune(body)[:LogRequestMaxSize])
		// }

		respLog := blw.body.String()
		// 打印响应日志
		if c.Writer.Status() == http.StatusOK {
			ngxLogInfo(c, respLog)
		} else {
			ngxLogError(c, respLog)
		}
	}
}

// ngx log
func ngxLogInfo(c *gin.Context, response string) {
	if logger == nil {
		return
	}

	list := listNgxField(c, zapcore.InfoLevel, response)
	logger.appLog.Info("", list...)
}

func ngxLogError(c *gin.Context, response string) {
	if logger == nil {
		return
	}

	list := listNgxField(c, zapcore.ErrorLevel, response)
	logger.appLog.Info("", list...)
}

func listNgxField(c *gin.Context, level zapcore.Level, response string) []zap.Field {
	var list []zap.Field
	//list = append(list, zap.String(LogKeyName, "ngx."+getLevelName(level)))
	//
	//list = append(list, zap.String(LogKeyTraceId, GetCtxParamString(c, LogKeyTraceId)))
	//list = append(list, zap.String(LogKeyRpcId, GetCtxParamString(c, LogKeyRpcId)))
	//
	//xFile, xLine := getFieldFileLine(1)
	//list = append(list, zap.String(LogKeyFile, xFile))
	//list = append(list, zap.String(LogKeyLine, xLine))
	//
	//list = append(list, zap.String(LogKeyExtra, GetCtxParamString(c, LogKeyExtra)))
	//
	//list = append(list, zap.String(LogKeyRequest, GetCtxParamString(c, LogKeyRequest)))
	//if response != "" {
	//	list = append(list, zap.String(LogKeyResponse, response))
	//}
	//
	//// 执行时间
	//duration := time.Since(GetCtxParamTime(c, CtxRequestStart)).Seconds()
	//list = append(list, zap.Float64(LogKeyDuration, duration))
	//
	//list = append(list, zap.String(LogKeyHeader, getCustomHeader(c)))

	return list
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
