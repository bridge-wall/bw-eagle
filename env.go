package eagle

import (
	"errors"
	"os"
)

const (
	DeployEnvDev  = "Development"
	DeployEnvTest = "Test"
	DeployEnvProd = "Production"

	NamespaceDev  = "dev"
	NamespaceTest = "test"
)

var (
	appEnv       string
	appMode      string
	appName      string
	appNamespace string
	appVersion   string
	appPodName   string
	appPodIp     string
)

// 检测运行环境
func CheckEnv(env string) error {
	appEnv = env

	appMode = os.Getenv("OCEAN_MODE")
	appName = os.Getenv("OCEAN_APP")
	appNamespace = os.Getenv("NAMESPACE")
	appVersion = os.Getenv("OCEAN_VERSION")
	appPodName = os.Getenv("HOSTNAME")
	appPodIp = os.Getenv("POD_IP")

	switch env {
		case DeployEnvProd:
			return nil
		case DeployEnvTest:
			return nil
		case DeployEnvDev:
			return nil
	}

	return errors.New("env error")
}

// 获取运行环境
func GetEnv() string {
	return appEnv
}

// 获取 OCEAN_MODE
func GetOceanMode() string {
	return appMode
}

func GetOceanApp() string {
	return appName
}

func GetOceanNamespace() string {
	return appNamespace
}

func GetOceanVersion() string {
	return appVersion
}

func GetOceanPodName() string {
	return appPodName
}

func GetOceanPodIp() string {
	return appPodIp
}
