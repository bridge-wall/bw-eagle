package eagle

import (
	"errors"
)

func ConfigNotFound(config string) error {
	return errors.New("[config] " + config + " not found")
}

func ConfigError(config string) error {
	return errors.New("[config] " + config + " error")
}

const (
	Success = 0

	ErrUnknown  = 99999 // 未知错误
	ErrInternal = 99998 // 内部错误
	ErrMysql    = 99997 // Mysql错误
	ErrRedis    = 99996 // Redis错误

	ErrParam = 10001 // 参数错误
)

// 错误码列表
var errCodes = make(map[int]string)

// 初始化错误码列表
func initErrMsg() {
	errCodes[Success] = "ok"

	errCodes[ErrUnknown] = "未知错误"
	errCodes[ErrInternal] = "内部错误"
	errCodes[ErrMysql] = "Mysql错误"
	errCodes[ErrRedis] = "Redis错误"

	errCodes[ErrParam] = "参数错误"

}

// 获取应用设置的错误信息
func ErrMsg(errCode int) string {
	v, ok := errCodes[errCode]
	if ok {
		return v
	}

	return errCodes[ErrUnknown]
}

// 设置应用的错误信息
func SetErrMsg(errCode int, message string) error {
	if _, ok := errCodes[errCode]; ok {
		panic("error code conflict")
	}

	errCodes[errCode] = message

	return nil
}
