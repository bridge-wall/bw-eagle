package eagle

// 路由注册
func registerRouter() {
	insHttpEngine.GET("/healthz", healthCheck)
}
