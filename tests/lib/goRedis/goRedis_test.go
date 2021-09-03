package gredis

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/gogf/gf/frame/g"
	"github.com/kdcer/go-lib/lib/goRedis"
)

func Test_redis(t *testing.T) {
	goRedis.InitRedis(&redis.Options{
		Addr:     g.Config().GetString("redis.addr"),
		Password: g.Config().GetString("redis.password"),
		DB:       0,
	})

	rdb := goRedis.Rdb
	rdb.Set("key1", "1", 0)
	rdb.Get("key1")
	rdb.Del("key1")

}
