package eagle

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type BwHttpEngine struct {
	*gin.Engine
}

var (
	insHttpEngine *gin.Engine
	httpHost      string
	httpPort      int
)

func NewHttp() error {
	if insHttpEngine != nil {
		return nil
	}

	httpHost = GetString("http.host")
	httpPort = GetInt("http.port")
	if httpPort == 0 {
		return ConfigNotFound("http.port")
	}

	gin.SetMode(gin.ReleaseMode)
	//if GetEnv() == DeployEnvDev {
	//	gin.SetMode(gin.DebugMode)
	//} else if GetEnv() == DeployEnvTest {
	//	gin.SetMode(gin.TestMode)
	//} else {
	//	gin.SetMode(gin.ReleaseMode)
	//}

	// 数字不要解析成float64
	binding.EnableDecoderUseNumber = true

	insHttpEngine = gin.New()
	HttpEngine = &BwHttpEngine{insHttpEngine}

	return nil
}

func RunHttp() error {
	if insHttpEngine == nil {
		return errors.New("http not init")
	}

	registerRouter()

	return insHttpEngine.Run(httpHost + ":" + strconv.Itoa(httpPort))
}

// 健康检查
type healthStatus struct {
	Status int    `json:"status"`
}

func healthCheck(c *gin.Context) {
	h := &healthStatus{
		Status: 200,
	}
	c.JSON(http.StatusOK, h)
}
