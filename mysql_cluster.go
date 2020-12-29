package eagle

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

type BwMysqlCluster map[int]*xorm.Engine

var (
	insMysqlCluster BwMysqlCluster
)

type mysqlClusterConfig struct {
	Id     int // 实例顺序，从0开始
	Config *mysqlConfig
}

// 简易版mysql集群多实例
// TODO 某实例宕机后如何让上层业务感知，摘除该实例
// TODO 某实例宕机一段时间后又恢复了如何让上层业务感知，启用该实例
func NewMysqlCluster(configs []*mysqlClusterConfig) (BwMysqlCluster, error) {
	if insMysqlCluster != nil {
		return insMysqlCluster, nil
	}

	insMysqlCluster = make(map[int]*xorm.Engine)

	for _, c := range configs {
		if _, ok := insMysqlCluster[c.Id]; ok {
			return nil, errors.New("mysql id conflict")
		}

		db, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", c.Config.User, c.Config.Password, c.Config.Host, c.Config.Port, c.Config.Database))
		if err != nil {
			return nil, err
		}
		err = db.Ping()
		if err != nil {
			return nil, err
		}

		// 设置空闲连接池中的最大连接数
		db.SetMaxIdleConns(c.Config.ConnMin)
		// 设置数据库连接最大打开数
		db.SetMaxOpenConns(c.Config.ConnMax)
		// 设置可重用连接的最长时间
		// 【注】一定要小于mysql服务端的保持超时时间，否则可能会被服务端关闭
		db.SetConnMaxLifetime(5 * time.Minute)

		db.SetMapper(core.SnakeMapper{})

		// 非生产环境开启sql日志
		if GetEnv() == DeployEnvDev || GetEnv() == DeployEnvTest {
			db.ShowSQL(true)
			db.ShowExecTime(true)
			// 测试环境连接保持超时较短
			if GetEnv() == DeployEnvTest {
				db.SetConnMaxLifetime(1 * time.Minute)
			}
		}

		insMysqlCluster[c.Id] = db
	}

	return insMysqlCluster, nil
}

// 返回所有实例
func (mc *BwMysqlCluster) Engines() map[int]*xorm.Engine {
	return insMysqlCluster
}

// 返回首个实例
func (mc *BwMysqlCluster) FirstEngine() *xorm.Engine {
	return insMysqlCluster[0]
}

// 返回集群中某一个实例，使用时最好判断是否为nil
func (mc *BwMysqlCluster) Engine(id int) *xorm.Engine {
	return insMysqlCluster[id]
}
