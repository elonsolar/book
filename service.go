package app

import (
	"fmt"
	"reflect"
)

// Service mgr all logical handlers
type Service struct {
	filters []func(name string, args []interface{})

	logicHandlerMap map[string]reflect.Value
	app             *App
}

func NewService(app *App) *Service {
	return &Service{
		app:             app,
		logicHandlerMap: make(map[string]reflect.Value, 0),
		filters: []func(name string, args []interface{}){
			func(name string, args []interface{}) {
				fmt.Println(name, args)
			},
		},
	}
}

func (s *Service) WithFilter(filter func(name string, args []interface{})) {

	s.filters = append(s.filters, filter)
}

// Register register a method with name
func (s *Service) Register(name string, fn reflect.Value) {

	if _, exist := s.logicHandlerMap[name]; exist {
		panic(fmt.Sprintf("method :%s already exist", name))
	}
	s.logicHandlerMap[name] = fn
}

// Call call func with name ,and execute
func (s *Service) Call(name string, data []interface{}) interface{} {

	for _, fn := range s.filters {
		fn(name, data)
	}

	var args []reflect.Value
	for _, arg := range data {
		args = append(args, reflect.ValueOf(arg))
	}
	fun, ok := s.logicHandlerMap[name]
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
