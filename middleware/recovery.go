package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

// 捕捉异常自动恢复，请求异常或堆栈状态信息写入日志
func HandlerRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				LogError(c, "[RecoveryPanic] [Request] "+fmt.Sprint(err))

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				request := strings.Replace(string(httpRequest), "\r", "|", -1)
				req := strings.Replace(request, "\n", "|", -1)
				LogError(c, "[RecoveryPanic] [Request] "+req)

				if stack {
					stack1 := strings.Replace(string(debug.Stack()), "\r", "|", -1)
					stack2 := strings.Replace(stack1, "\n", "|", -1)
					LogError(c, "[RecoveryPanic] [Request] "+stack2)
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
