package douyin

import (
	"github.com/fastwego/microapp"
	"github.com/gogf/gf/frame/g"
	"github.com/kdcer/go-lib/lib/douying"
	"testing"
)

//初始化 抖音小程序配置
func Test_douyin(t *testing.T) {
	douying.InitMicroapp(microapp.Config{
		AppId:     g.Config().GetString("douyin.AppId"),
		AppSecret: g.Config().GetString("douyin.AppSecret"),
	})
}
