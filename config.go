package eagle

import (
	"time"

	"github.com/spf13/viper"
)

var (
	insConfig *viper.Viper
)

func NewConfig() error {
	if insConfig != nil {
		return nil
	}

	insConfig = viper.New()
	insConfig.SetConfigName("app")
	insConfig.AddConfigPath("./conf")
	err := insConfig.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func GetConfig(key string) interface{} {
	return insConfig.Get(key)
}

func GetString(key string) string {
	return insConfig.GetString(key)
}

func GetBool(key string) bool {
	return insConfig.GetBool(key)
}

func GetInt(key string) int {
	return insConfig.GetInt(key)
}

func GetInt32(key string) int32 {
	return insConfig.GetInt32(key)
}

func GetInt64(key string) int64 {
	return insConfig.GetInt64(key)
}

func GetUint(key string) uint {
	return insConfig.GetUint(key)
}

func GetUint32(key string) uint32 {
	return insConfig.GetUint32(key)
}

func GetUint64(key string) uint64 {
	return insConfig.GetUint64(key)
}

func GetFloat64(key string) float64 {
	return insConfig.GetFloat64(key)
}

func GetTime(key string) time.Time {
	return insConfig.GetTime(key)
}

func GetDuration(key string) time.Duration {
	return insConfig.GetDuration(key)
}

// func GetIntSlice(key string) []int {
// 	return insConfig.GetIntSlice(key)
// }

func GetStringSlice(key string) []string {
	return insConfig.GetStringSlice(key)
}

func GetStringMap(key string) map[string]interface{} {
	return insConfig.GetStringMap(key)
}

func GetStringMapString(key string) map[string]string {
	return insConfig.GetStringMapString(key)
}

func GetStringMapStringSlice(key string) map[string][]string {
	return insConfig.GetStringMapStringSlice(key)
}

// func GetSizeInBytes(key string) uint {
// 	return insConfig.GetSizeInBytes(key)
// }
