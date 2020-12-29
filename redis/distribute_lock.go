package redis

import (
	"errors"
)

// 分布式锁
// TODO 先简单实现，后续再优化

type DistributeLock struct {
	expireSecond int
}

type DistributeLockOpt func(*DistributeLock)

// 初始化
func NewDistributeLock(opts ...DistributeLockOpt) *DistributeLock {
	d := &DistributeLock{}

	for _, o := range opts {
		o(d)
	}

	return d
}

// 设置过期时间秒数
func SetExpireSecond(expire int) DistributeLockOpt {
	return func(d *DistributeLock) {
		d.expireSecond = expire
	}
}

// 加锁
func (d *DistributeLock) Lock(key, value string) (bool, error) {
	if d.expireSecond > 0 {
		reply, err := insRedis.String("SET", key, value, "EX", d.expireSecond, "NX")
		if err != nil {
			return false, err
		}

		return reply == "OK", nil
	}

	return false, errors.New("expire second not set")
}

// 解锁
// TODO 验证加锁和解锁是同一个
func (d *DistributeLock) UnLock(key, value string) (int, error) {
	return insRedis.Int("DEL", key)
}
