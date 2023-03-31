package basis

var apiInitializerRegister *InitializeRegister = new(InitializeRegister)

func RegisterApi(ai Initializer) {
	apiInitializerRegister.Register(ai)
}

func GetApiInitializers() []Initializer {
	return apiInitializerRegister.Initializers
}

type WebApiStarter struct {
	BaseStarter
}

func (w *WebApiStarter) Setup(ctx StarterContext) {
	for _, v := range GetApiInitializers() {
		v.Init()
	}
}
