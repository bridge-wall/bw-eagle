package eagle

import (
	"flag"
	"github.com/rhonin-cd/rhonin-eagle/middleware"
	"github.com/rhonin-cd/rhonin-eagle/redis"
	"time"
)

func New() error {
	flag.Parse()

	err := NewConfig()
	if err != nil {
		panic(err)
	}

	initErrMsg()

	httpConfig := GetConfig("http")
	if httpConfig != nil {
		err = NewHttp()
		if err != nil {
			panic(err)
		}
	}
	// 日志配置
	logConfig := GetConfig("log")
	if logConfig != nil {
		err = loadLog()
		if err != nil {
			panic(err)
		}
	}

	// 单实例mysql
	mysqlConfig := GetConfig("mysql")
	if mysqlConfig != nil {
		err = loadMysql()
		if err != nil {
			panic(err)
		}
	}

	// 单实例redis启用
	redisConfig := GetConfig("redis")
	if redisConfig != nil {
		err = loadRedis()
		if err != nil {
			panic(err)
		}
	}

	// 多实例mysql
	mcConfig := GetStringMap("mysql-cluster")
	if mcConfig != nil {
		err = loadMysqlCluster(mcConfig)
		if err != nil {
			panic(err)
		}
	}

	loadMiddleware()
	return nil
}
// 启动服务
func Run() error {
	err := RunHttp()
	if err != nil {
		panic(err)
	}

	return nil
}

func loadLog() error {
	config := &middleware.LogConfig{}

	config.AppPath = GetString("log.app")
	if config.AppPath == "" {
		return ConfigNotFound("log.app")
	}

	config.Extra.AppName = GetOceanApp()
	config.Extra.Env = GetEnv()
	config.Extra.Mode = GetOceanMode()
	config.Extra.Version = GetOceanVersion()
	config.Extra.Namespace = GetOceanNamespace()
	config.Extra.PodName = GetOceanPodName()
	config.Extra.PodIp = GetOceanPodIp()

	_, err := middleware.NewLog(config)
	if err != nil {
		return err
	}

	return nil
}

func loadMysql() error {
	config := &mysqlConfig{}

	config.Host = GetString("mysql.host")
	if config.Host == "" {
		return ConfigNotFound("mysql.host")
	}

	config.Port = GetInt("mysql.port")
	if config.Port == 0 {
		return ConfigNotFound("mysql.port")
	}

	config.User = GetString("mysql.user")
	if config.User == "" {
		return ConfigNotFound("mysql.user")
	}

	config.Password = GetString("mysql.password")
	if config.Password == "" {
		return ConfigNotFound("mysql.password")
	}

	config.Database = GetString("mysql.database")
	if config.Database == "" {
		return ConfigNotFound("mysql.database")
	}

	config.ConnMin = GetInt("mysql.conn_min")
	if config.ConnMin == 0 {
		return ConfigNotFound("mysql.conn_min")
	}

	config.ConnMax = GetInt("mysql.conn_max")
	if config.ConnMax == 0 {
		return ConfigNotFound("mysql.conn_max")
	}

	if config.ConnMax < config.ConnMin {
		return ConfigError("mysql.conn_min, mysql.conn_max")
	}

	var err error
	MysqlEngine, err = NewMysql(config)
	if err != nil {
		return err
	}

	// 非生产环境开启sql日志
	if GetEnv() == DeployEnvDev || GetEnv() == DeployEnvTest {
		insMysql.ShowSQL(true)
		insMysql.ShowExecTime(true)
	}

	return nil
}

func loadRedis() error {
	config := &redis.Config{}

	config.Address = GetString("redis.address")
	if config.Address == "" {
		return ConfigNotFound("redis.address")
	}

	config.Auth = GetString("redis.auth")
	config.Database = GetInt("redis.database")

	config.MaxIdle = GetInt("redis.conn_min")
	if config.MaxIdle == 0 {
		return ConfigNotFound("redis.conn_min")
	}
	config.MaxActive = GetInt("redis.conn_max")

	config.IdleTimeout = 30 * time.Second
	config.ConnectTimeout = 2 * time.Second
	config.ReadTimeout = 2 * time.Second
	config.WriteTimeout = 2 * time.Second

	var err error
	RedisEngine, err = redis.NewRedis(config)
	if err != nil {
		return err
	}

	return nil
}

// mysql集群配置
func loadMysqlCluster(clusterConfig map[string]interface{}) error {
	var configs []*mysqlClusterConfig
	for _, v := range clusterConfig {
		e := v.(map[string]interface{})
		if e["id"] == nil {
			return ConfigNotFound("mysql.id")
		}
		if e["host"] == nil {
			return ConfigNotFound("mysql.host")
		}
		if e["port"] == nil {
			return ConfigNotFound("mysql.port")
		}
		if e["user"] == nil {
			return ConfigNotFound("mysql.user")
		}
		if e["password"] == nil {
			return ConfigNotFound("mysql.password")
		}
		if e["database"] == nil {
			return ConfigNotFound("mysql.database")
		}
		if e["conn_min"] == nil {
			return ConfigNotFound("mysql.conn_min")
		}
		if e["conn_max"] == nil {
			return ConfigNotFound("mysql.conn_max")
		}

		conf := &mysqlConfig{}
		conf.Host = e["host"].(string)
		conf.Port = int(e["port"].(int64))
		conf.User = e["user"].(string)
		conf.Password = e["password"].(string)
		conf.Database = e["database"].(string)
		conf.ConnMin = int(e["conn_min"].(int64))
		conf.ConnMax = int(e["conn_max"].(int64))
		config := &mysqlClusterConfig{
			Id:     int(e["id"].(int64)),
			Config: conf,
		}
		configs = append(configs, config)
	}

	var err error
	MysqlCluster, err = NewMysqlCluster(configs)
	if err != nil {
		return err
	}

	return nil
}