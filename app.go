package app

import (
	"fmt"
	"reflect"
)

type App struct {
	cfg        *Config
	Controller *Controller
	Service    *Service
	Dao        *Dao
	errors     []error

	filters []func(name string, args []interface{})

	handlerMap map[string]reflect.Value
}

func NewApp(cfg *Config) *App {

	app := &App{
		cfg:        cfg,
		handlerMap: make(map[string]reflect.Value, 0),
		filters: []func(name string, args []interface{}){
			func(name string, args []interface{}) {
				fmt.Printf("------------拦截所有 方法 和参数--------- 方法名:%s \n,  参数：%v", name, args)
			},
		},
	}
	if cfg.ControllerCfg != nil {

		ctl := newController(cfg.ControllerCfg, app)
		app.Controller = ctl
	}

	if cfg.DaoCfg != nil {

		dao := newDao(cfg.DaoCfg, app)
		app.Dao = dao

	}
	return app
}

func (s *App) WithFilter(filter func(name string, args []interface{})) {

	s.filters = append(s.filters, filter)
}

// Register register a method with name
func (s *App) Register(name string, fn reflect.Value) {

	if _, exist := s.handlerMap[name]; exist {
		panic(fmt.Sprintf("method :%s already exist", name))
	}
	s.handlerMap[name] = fn
}

// Call call func with name ,and execute
func (s *App) Call(name string, data []interface{}) interface{} {

	for _, fn := range s.filters {
		fn(name, data)
	}

	var args []reflect.Value
	for _, arg := range data {
		args = append(args, reflect.ValueOf(arg))
	}
	fun, ok := s.handlerMap[name]
	if !ok {
		panic(fmt.Sprintf("no such method name :%s", name))
	}

	ret := fun.Call(args)
	var result []interface{}

	for _, r := range ret {
		result = append(result, r.Interface())
	}
	return result
}

func (a *App) Start() {

	a.Controller.Start()
}

func (a *App) Error() {

	fmt.Println(a.errors)
}
