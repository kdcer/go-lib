package goRedis

import (
	"github.com/go-redis/redis"
)

var Rdb *redis.Client

func InitRedis(options *redis.Options) {
	Rdb = redis.NewClient(options)
}

func CheckAndDel(key, value string) (int, error) {
	cmd := Rdb.Eval(`if redis.call("get",KEYS[1]) == ARGV[1]
										then
											return redis.call("del",KEYS[1])
										else
											return 0
										end`, []string{key}, value)
	res, e := cmd.Int()
	return res, e
}
