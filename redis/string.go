package redis

// 字符串命令
func Set(key, value string) (string, error) {
	return insRedis.String("SET", key, value)
}

func Get(key string) (string, error) {
	return insRedis.String("GET", key)
}

func SetEx(key, value string, seconds int) (string, error) {
	return insRedis.String("SETEX", key, seconds, value)
}
