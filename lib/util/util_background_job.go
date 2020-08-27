package util

import (
	"go-lib/lib/gpool"

	"github.com/gogf/gf/os/glog"
)

//协程池
var pool *gpool.Pool

var defaultPoolLimitSize uint64 = 30

func InitPool(poolLimitSize uint64) {

	_pool, err := gpool.NewPool(poolLimitSize)
	if err != nil {
		glog.Error("初始化gpool.Pool失败 poolLimitSize=", poolLimitSize)
		panic(err)
	}

	pool = _pool
	glog.Info("gpool.Pool===>>> cap=", pool.GetCap())
}

//获取gpools
func GetGpool() *gpool.Pool {
	if pool == nil {
		glog.Info("gpool.Pool未初始化, 默认初始化poolLimitSize=", defaultPoolLimitSize)
		InitPool(defaultPoolLimitSize)
	}
	return pool
}

// 添加后台任务
func AddTask(task *gpool.Task) {
	gpool.GetDefaultPool().Put(task)
}

//添加后台任务
func AddTaskByFuc(handler func(args ...interface{}), args ...interface{}) {
	gpool.GetDefaultPool().Put(&gpool.Task{
		Handler: handler,
		Params:  args,
	})
}

// 关闭协程池
func ClosePool(task *gpool.Task) {
	gpool.ClosePool()
}
