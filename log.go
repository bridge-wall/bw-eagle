package eagle

import (
	"context"

	"github.com/bridge-wall/bw-eagle/middleware"
)

func LogDebug(ctx context.Context, a interface{}) {
	// 生产环境不能打印Debug日志
	//if GetEnv() == DeployEnvProd {
	//	return
	//}

	middleware.LogDebug(ctx, a)
}

func LogInfo(ctx context.Context, a interface{}) {
	middleware.LogInfo(ctx, a)
}

func LogWarn(ctx context.Context, a interface{}) {
	middleware.LogWarn(ctx, a)
}

func LogError(ctx context.Context, a interface{}) {
	middleware.LogError(ctx, a)
}

func LogCritical(ctx context.Context, a interface{}) {
	middleware.LogCritical(ctx, a)
}

// 设置日志打印哪些header，非空时打印到日志中
func LogAddCustomHeader(header string) {
	middleware.LogAddCustomHeader(header)
}
