package eagle

import (
	"github.com/blust/bw-eagle/middleware"
	"github.com/gin-contrib/pprof"
)


// 加载中间件
func loadMiddleware() {
	// 开启异常恢复，并写入日志
	// 开发模式不用开启，方便查找问题
	if GetEnv() != DeployEnvDev {
		// 默认开启
		recovery := GetConfig("core.recovery")
		if recovery == nil || recovery.(int64) == 1 {
			insHttpEngine.Use(middleware.HandlerRecovery(true))
		}
	}

	prof := GetInt("http.pprof")
	if prof == 1 {
		pprof.Register(insHttpEngine)
	}

	// 上下文初始化
	ctxInit := GetConfig("core.context_init")
	if ctxInit == nil || ctxInit.(int64) == 1 {
		insHttpEngine.Use(middleware.HandlerInitContext())
	}
}
