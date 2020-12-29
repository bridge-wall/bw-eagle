package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

// 浮点数四舍五入，保留N位小数
func Round(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}

// 生成随机数
func GetRandomNum(min, max int) int {
	return rand.Intn(max) + min
}

// 生成随机数
func GetRandomUuid() string {
	ul, err := uuid.NewV4()
	if err != nil {
		return strconv.FormatInt(time.Now().Unix(), 10)
	}

	return Md5(ul.String() + strconv.Itoa(GetRandomNum(1, 9999999)))
}
