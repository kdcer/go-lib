package douying

import (
	"github.com/fastwego/microapp"
)

var MicroappApp *microapp.MicroApp

func InitMicroapp(conf microapp.Config) {
	MicroappApp = microapp.New(conf)
}
