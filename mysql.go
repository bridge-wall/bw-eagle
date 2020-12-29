package eagle

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

type MysqlClient struct {
	*xorm.Engine
}

type mysqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string

	ConnMin int
	ConnMax int
}

var (
	insMysql *MysqlClient
)

func NewMysql(config *mysqlConfig) (*MysqlClient, error) {
	if insMysql != nil {
		return insMysql, nil
	}

	db, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.User, config.Password, config.Host, config.Port, config.Database))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	insMysql = &MysqlClient{db}

	// 设置空闲连接池中的最大连接数
	insMysql.SetMaxIdleConns(config.ConnMin)
	// 设置数据库连接最大打开数
	insMysql.SetMaxOpenConns(config.ConnMax)
	// 设置可重用连接的最长时间
	// 【注】一定要小于mysql服务端的保持超时时间，否则可能会被服务端关闭
	insMysql.SetConnMaxLifetime(5 * time.Minute)

	insMysql.SetMapper(core.SnakeMapper{})

	return insMysql, nil
}

type MysqlStatus struct {
	// 数据库状态统计
	DDStats *sql.DBStats

	// TODO 时间统计

}

func (e *MysqlClient) Status() *MysqlStatus {
	dbStats := insMysql.DB().Stats()

	return &MysqlStatus{
		DDStats: &dbStats,
	}
}

type Model struct {
	table string
}

func NewModel(table string) *Model {
	return &Model{
		table: table,
	}
}

func (m *Model) Insert(data interface{}) (int64, error) {
	return insMysql.Table(m.table).Insert(data)
}

func (m *Model) InsertOne(data interface{}) (int64, error) {
	return insMysql.Table(m.table).InsertOne(data)
}

func (m *Model) InsertMulti(rowsSlicePtr interface{}) (int64, error) {
	return insMysql.Table(m.table).InsertMulti(rowsSlicePtr)
}

func (m *Model) Update(cond, fields map[string]interface{}) (int64, error) {
	affected, err := insMysql.Table(m.table).Update(fields, cond)
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func (m *Model) Get(cond map[string]interface{}, result interface{}) (err error) {
	var has bool
	has, err = insMysql.Table(m.table).Where(cond).Get(result)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	return nil
}

func (m *Model) Count(cond map[string]interface{}) (count int64, err error) {
	if len(cond) > 0 {
		return insMysql.Table(m.table).Where(cond).Count()
	} else {
		return insMysql.Table(m.table).Count()
	}
}

func (m *Model) Find(record interface{}, cond ...interface{}) error {
	return insMysql.Table(m.table).Find(record, cond)
}

func (m *Model) FindAndCount(record interface{}, cond ...interface{}) (int64, error) {
	return insMysql.Table(m.table).FindAndCount(record, cond)
}
