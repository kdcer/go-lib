package alone

import (
	"github.com/gomodule/redigo/redis"
	"github.com/kdcer/go-lib/lib/gredis/mode"
)

type options struct {
	addr     string
	poolOpts []mode.PoolOption
	dialOpts []redis.DialOption
}

type OptFunc func(opts *options)

func Addr(value string) OptFunc {
	return func(opts *options) {
		opts.addr = value
	}
}

func PoolOpts(value ...mode.PoolOption) OptFunc {
	return func(opts *options) {
		for _, poolOpt := range value {
			opts.poolOpts = append(opts.poolOpts, poolOpt)
		}
	}
}

func DialOpts(value ...redis.DialOption) OptFunc {
	return func(opts *options) {
		for _, dialOpt := range value {
			opts.dialOpts = append(opts.dialOpts, dialOpt)
		}
	}
}
