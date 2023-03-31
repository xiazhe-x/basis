package base

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"github.com/xiazhe-x/basis"
)

var (
	props      kvs.ConfigSource
	SystemConf *systemConf
)

func Props() kvs.ConfigSource {
	Check(props)
	return props
}

type PropsStarter struct {
	basis.BaseStarter
}

type systemConf struct {
	Port        int
	Name        string
	TokenSecret string
	MqttAddr    string
	MqttUser    string
	MqttPwd     string
	Environment string
}

func (p *PropsStarter) Init(ctx basis.StarterContext) {
	props = ctx.Props()
	SystemConf = new(systemConf)
	err := kvs.Unmarshal(Props(), SystemConf, "app")
	if err != nil {
		panic(err)
	}
	log.Info("初始化配置.")
}
