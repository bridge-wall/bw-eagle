package eagle

import (
	"context"
	"github.com/bridge-wall/bw-eagle/redis"

	"github.com/bridge-wall/bw-eagle/middleware"
	"github.com/gin-gonic/gin"
)

var (
	// Http
	HttpEngine *BwHttpEngine
	// Mysql
	MysqlEngine *MysqlClient
	// Redis
	RedisEngine *redis.BwRedisClient
	// MysqlCluster
	MysqlCluster BwMysqlCluster
)

// 接口开启访问日志
func AccessLogHandler() gin.HandlerFunc {
	return middleware.HandlerAccessLog()
}


// 获取请求URI
func GetRequestURI(ctx context.Context) string {
	if ctx != nil {
		if c, ok := ctx.(*gin.Context); ok {
			return c.Request.RequestURI
		}
	}

	return ""
}

// 获取请求body
func GetRequestBody(ctx context.Context) string {
	if ctx != nil {
		if c, ok := ctx.(*gin.Context); ok {
			return middleware.GetCtxParamString(c, middleware.CtxRequestBody)
		}
	}

	return ""
}