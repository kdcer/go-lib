package goRedis

import (
	"github.com/go-redis/redis"
	"github.com/gogf/gf/frame/g"
)

var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     g.Config().GetString("redis.addr"),
		Password: g.Config().GetString("redis.password"), // no password set
		DB:       0,                                      // use default DB
	})
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
