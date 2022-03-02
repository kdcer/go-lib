package casbin

import (
	"strings"
	"sync"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

var (
	Enforcer *casbin.Enforcer
	syncOnce sync.Once
)

func New(driverName, dataSourceName, confPath string) {
	syncOnce.Do(func() {
		// 要使用自己定义的数据库rbac_db,最后的true很重要.默认为false,使用缺省的数据库名casbin,不存在则创建
		a, err := xormadapter.NewAdapter(driverName, dataSourceName, true)
		if err != nil {
			glog.Error("casbin连接数据库错误: %v", err)
			panic(err)
		}
		e, err := casbin.NewEnforcer(confPath, a)
		if err != nil {
			glog.Error("初始化casbin错误: %v", err)
			panic(err)
		}
		Enforcer = e
	})
}

// New2 goframe配置文件专用
func New2(confPath string) {
	link := g.Config().GetString("database.link")
	if len(link) == 0 {
		panic("casbin数据库连接为空")
	}
	links := strings.Split(link, ":")
	if len(links) == 0 {
		panic("casbin数据库连接错误")
	}
	driverName := links[0]
	dataSourceName := strings.Replace(link, driverName+":", "", 1)
	New(driverName, dataSourceName, confPath)
}
