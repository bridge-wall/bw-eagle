package redis

// 集合命令

// 向集合添加一个或多个成员
func SAdd(key string, members []interface{}) (int, error) {
	args := make([]interface{}, 0, len(members)+1)
	args = append(args, key)
	args = append(args, members...)

	return insRedis.Int("SADD", args...)
}

// 获取集合的成员数量
func SCard(key string) (int, error) {
	return insRedis.Int("SCARD", key)
}

// 返回集合中的所有成员
func SMembers(key string) ([]string, error) {
	result, err := insRedis.Strings("SMEMBERS", key)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 判断元素是否是集合的成员
func SIsMember(key, member string) (int, error) {
	return insRedis.Int("SISMEMBER", key, member)
}

// 返回给定所有集合的交集
func SInter(keys []interface{}) ([]string, error) {
	return insRedis.Strings("SINTER", keys...)
}

// 移除集合中一个或多个成员
func SRem(key string, members []interface{}) (int, error) {
	args := make([]interface{}, 0, len(members)+1)
	args = append(args, key)
	args = append(args, members...)

	return insRedis.Int("SREM", args...)
}
