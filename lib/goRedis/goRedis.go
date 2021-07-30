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
