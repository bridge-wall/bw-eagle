package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type BwRedisClient struct {
	*redis.Pool
}

type Config struct {
	Address  string
	Auth     string
	Database int

	MaxIdle        int
	MaxActive      int
	IdleTimeout    time.Duration
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

var (
	insRedis *BwRedisClient
)

func NewRedis(config *Config) (*BwRedisClient, error) {
	if insRedis != nil {
		return insRedis, nil
	}

	pool := &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: config.IdleTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Address,
				redis.DialConnectTimeout(config.ConnectTimeout),
				redis.DialReadTimeout(config.ReadTimeout),
				redis.DialWriteTimeout(config.WriteTimeout),
			)
			if err != nil {
				return nil, err
			}
			if config.Auth != "" {
				_, err = c.Do("AUTH", config.Auth)
				if err != nil {
					c.Close()
					return nil, err
				}
			}

			if config.Database > 0 {
				_, err = c.Do("SELECT", config.Database)
				if err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		// TestOnBorrow: func(c redis.Conn, t time.Time) error {
		// 	if time.Since(t) < time.Minute {
		// 		return nil
		// 	}
		// 	_, err := c.Do("PING")
		// 	return err
		// },
	}

	insRedis = &BwRedisClient{
		pool,
	}

	return insRedis, nil
}

// func CloseRedis() {
// 	NvRedis.pool.Close()
// }

type NvRedisStatus struct {
	// 连接池状态统计
	PoolActiveCount int
	PoolIdleCount   int

	// TODO 时间统计

}

func (r *BwRedisClient) Status() *NvRedisStatus {
	poolStats := insRedis.Stats()

	return &NvRedisStatus{
		PoolActiveCount: poolStats.ActiveCount,
		PoolIdleCount:   poolStats.IdleCount,
	}
}

func (r *BwRedisClient) ActiveCount() int {
	return insRedis.ActiveCount()
}

func (r *BwRedisClient) IdleCount() int {
	return insRedis.IdleCount()
}

// 返回 int
func (r *BwRedisClient) Int(cmd string, args ...interface{}) (int, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Int(reply, e)
	if err == redis.ErrNil {
		return 0, nil
	}

	return v, err
}

// 返回 int64
func (r *BwRedisClient) Int64(cmd string, args ...interface{}) (int64, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Int64(reply, e)
	if err == redis.ErrNil {
		return 0, nil
	}

	return v, err
}

// 返回 uint64
func (r *BwRedisClient) Uint64(cmd string, args ...interface{}) (uint64, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Uint64(reply, e)
	if err == redis.ErrNil {
		return 0, nil
	}

	return v, err
}

// 返回 float64
func (r *BwRedisClient) Float64(cmd string, args ...interface{}) (float64, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Float64(reply, e)
	if err == redis.ErrNil {
		return 0, nil
	}

	return v, err
}

// 返回 string
func (r *BwRedisClient) String(cmd string, args ...interface{}) (string, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.String(reply, e)
	if err == redis.ErrNil {
		return "", nil
	}

	return v, err
}

// 返回 bytes
func (r *BwRedisClient) Bytes(cmd string, args ...interface{}) ([]byte, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Bytes(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

// 返回 bool
func (r *BwRedisClient) Bool(cmd string, args ...interface{}) (bool, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Bool(reply, e)
	if err == redis.ErrNil {
		return false, nil
	}

	return v, err
}

// 返回 []interface{}
func (r *BwRedisClient) Values(cmd string, args ...interface{}) ([]interface{}, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Values(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

// 返回 []float64
func (r *BwRedisClient) Float64s(cmd string, args ...interface{}) ([]float64, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Float64s(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

// 返回 []string
func (r *BwRedisClient) Strings(cmd string, args ...interface{}) ([]string, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Strings(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

// 返回 [][]byte
func (r *BwRedisClient) ByteSlices(cmd string, args ...interface{}) ([][]byte, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.ByteSlices(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

// 返回 []int64
func (r *BwRedisClient) Int64s(cmd string, args ...interface{}) ([]int64, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Int64s(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

// 返回 []int
func (r *BwRedisClient) Ints(cmd string, args ...interface{}) ([]int, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Ints(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

// 返回 map[string]string
func (r *BwRedisClient) StringMap(cmd string, args ...interface{}) (map[string]string, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.StringMap(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

// 返回 map[string]int
func (r *BwRedisClient) IntMap(cmd string, args ...interface{}) (map[string]int, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.IntMap(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

// 返回 map[string]int64
func (r *BwRedisClient) Int64Map(cmd string, args ...interface{}) (map[string]int64, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Int64Map(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

// 返回 positions
func (r *BwRedisClient) Positions(cmd string, args ...interface{}) ([]*[2]float64, error) {
	conn := r.Get()
	defer conn.Close()

	reply, e := conn.Do(cmd, args...)
	v, err := redis.Positions(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

//
// 短连接
type NvRedisShortClient struct {
	redis.Conn
}

var (
	insRedisShort *NvRedisShortClient
)

func NewRedisShort(address string) (*NvRedisShortClient, error) {

	c, err := redis.DialTimeout("tcp", address, 1*time.Second, 1*time.Second, 1*time.Second)
	if err != nil {
		return nil, err
	}

	insRedisShort = &NvRedisShortClient{c}

	return insRedisShort, nil
}

func (r *NvRedisShortClient) CloseConn() error {
	return r.Close()
}

func (r *NvRedisShortClient) ReplyInt(cmd string, args ...interface{}) (int, error) {
	reply, e := r.Do(cmd, args...)
	v, err := redis.Int(reply, e)
	if err == redis.ErrNil {
		return 0, nil
	}

	return v, err
}

func (r *NvRedisShortClient) ReplyInt64(cmd string, args ...interface{}) (int64, error) {
	reply, e := r.Do(cmd, args...)
	v, err := redis.Int64(reply, e)
	if err == redis.ErrNil {
		return 0, nil
	}

	return v, err
}

func (r *NvRedisShortClient) ReplyUint64(cmd string, args ...interface{}) (uint64, error) {
	reply, e := r.Do(cmd, args...)
	v, err := redis.Uint64(reply, e)
	if err == redis.ErrNil {
		return 0, nil
	}

	return v, err
}

func (r *NvRedisShortClient) ReplyFloat64(cmd string, args ...interface{}) (float64, error) {
	reply, e := r.Do(cmd, args...)
	v, err := redis.Float64(reply, e)
	if err == redis.ErrNil {
		return 0, nil
	}

	return v, err
}

func (r *NvRedisShortClient) ReplyString(cmd string, args ...interface{}) (string, error) {
	reply, e := r.Do(cmd, args...)
	v, err := redis.String(reply, e)
	if err == redis.ErrNil {
		return "", nil
	}

	return v, err
}

func (r *NvRedisShortClient) ReplyBytes(cmd string, args ...interface{}) ([]byte, error) {
	reply, e := r.Do(cmd, args...)
	v, err := redis.Bytes(reply, e)
	if err == redis.ErrNil {
		return nil, nil
	}

	return v, err
}

func (r *NvRedisShortClient) ReplyBool(cmd string, args ...interface{}) (bool, error) {
	reply, e := r.Do(cmd, args...)
	v, err := redis.Bool(reply, e)
	if err == redis.ErrNil {
		return false, nil
	}

	return v, err
}
