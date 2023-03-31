package basis

import (
	"github.com/tietang/props/kvs"
)

type BootApplication struct {
	IsTest     bool
	conf       kvs.ConfigSource
	starterCtx StarterContext
}

func New(conf kvs.ConfigSource) *BootApplication {
	e := &BootApplication{conf: conf, starterCtx: StarterContext{}}
	e.starterCtx.SetProps(conf)
	return e
}

func (b *BootApplication) Start() {
	b.init()
	b.setup()
	b.start()
}

func (e *BootApplication) init() {

	for _, v := range GetStarters() {
		// typ := reflect.TypeOf(v)
		// logrus.Debugf("Initializing: PriorityGroup=%d,Priority=%d,type=%s", v.PriorityGroup(), v.Priority(), typ.String())
		v.Init(e.starterCtx)
	}
}

func (e *BootApplication) setup() {

	//log.Info("Setup starters...")
	for _, v := range GetStarters() {
		//typ := reflect.TypeOf(v)
		//log.Debug("Setup: ", typ.String())
		v.Setup(e.starterCtx)
	}

}

func (e *BootApplication) start() {

	//log.Info("Starting starters...")
	for i, v := range GetStarters() {

		//typ := reflect.TypeOf(v)
		//log.Debug("Starting: ", typ.String())
		if v.StartBlocking() {
			if i+1 == len(GetStarters()) {
				v.Start(e.starterCtx)
			} else {
				go v.Start(e.starterCtx)
			}
		} else {
			v.Start(e.starterCtx)
		}

	}
}

func (e *BootApplication) Stop() {

	//log.Info("Stoping starters...")
	for _, v := range GetStarters() {
		//typ := reflect.TypeOf(v)
		//log.Debug("Stoping: ", typ.String())
		v.Stop(e.starterCtx)
	}
}
