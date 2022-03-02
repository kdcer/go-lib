package casbin

import (
	"fmt"
	"testing"

	"github.com/kdcer/go-lib/lib/casbin"
)

func Test_Casbin(t *testing.T) {
	casbin.New("mysql", "root:123456@tcp(127.0.0.1:3306)/goblog?charset=utf8", "./rbac_models.conf")

	e := casbin.Enforcer
	list := e.GetPolicy()
	fmt.Println(list)
}
