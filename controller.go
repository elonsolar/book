package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControllerConfig struct {
	Port int
}

// Handler gin router handlers
type Handler struct {
	Method string
	Url    string
	fn     func(app *App) gin.HandlerFunc
	Name   string
}

// MidleWare gin midleware handlers
type MidleWare struct {
	Name   string
	fn     gin.HandlerFunc
	before string
	after  string
	c      *Controller
}

func (md *MidleWare) Register(name string, fn gin.HandlerFunc) error {
	md.Name = name
	md.fn = fn

	return md.c.SortMiddleWare()
}

type Controller struct {
	cfg *ControllerConfig

	*gin.Engine

	// all of the midlWare
	midleWares []*MidleWare
	// ignore midleware
	IgnoredMidleWare map[string]bool

	// handlers
	handlers []*Handler

	publicHandlers []*Handler
	app            *App
}

func NewController(cfg *ControllerConfig, app *App) *Controller {
	return &Controller{
		cfg:        cfg,
		Engine:     gin.Default(),
		midleWares: []*MidleWare{}, //default
		// {Name: "CORS", fn: nil},
		// {Name: "JWT", fn: nil},
		// {Name: "REPEATCHECK", fn: nil}},
		IgnoredMidleWare: make(map[string]bool),
		app:              app,
	}
}

// UseMidleWare     ["JWT","CORS"]
// use/unuse  the given midleware
func (r *Controller) UnUseMidleWare(middlewares []string) *Controller {

	for _, midle := range middlewares {

		r.IgnoredMidleWare[midle] = false
	}
	return r
}

// RegisterMidleWare
func (r *Controller) RegisterMidleWare(name string, fn gin.HandlerFunc) error {

	return (&MidleWare{}).Register(name, fn)
}

func (r *Controller) AfterMidleWare(name string) *MidleWare {
	return &MidleWare{after: name, c: r}
}

func (r *Controller) BeforeMidleWare(name string) *MidleWare {
	return &MidleWare{before: name, c: r}
}

func (r *Controller) SortMiddleWare() error {

	midleWares := r.midleWares

	var getIndex = func(tname string) int {

		for i, mid := range midleWares {
			if mid.Name == tname {
				return i
			}
		}
		return -1
	}

	for _, mid := range r.midleWares {

		if mid.before != "" {

			if index := getIndex(mid.before); index != -1 {
				midleWares = append(midleWares[:index], append([]*MidleWare{mid}, midleWares[index:]...)...)
			} else {
				return fmt.Errorf("未找到%s", mid.before)
			}
		}

		if mid.after != "" {

			if index := getIndex(mid.before); index != -1 {
				if index >= len(midleWares)-1 {
					midleWares = append(midleWares, mid)
				} else {
					midleWares = append(midleWares[:index+1], append([]*MidleWare{mid}, midleWares[index+1:]...)...)
				}
			} else {

				return fmt.Errorf("未找到%s", mid.after)
			}
		}
	}

	return nil
}

// WithHandlers
func (r *Controller) WithHandlers(handlers []*Handler, public bool) *Controller {

	if public {
		r.publicHandlers = append(r.publicHandlers, handlers...)
	}

	r.handlers = append(r.handlers, handlers...)
	return r
}

func (r *Controller) initRouter(handlers []*Handler) {

	for _, handler := range handlers {

		switch handler.Method {

		case http.MethodPost:

			r.POST(handler.Url, handler.fn(r.app)) // 处理器注入 app,可以调用server 函数
		case http.MethodGet:

			r.GET(handler.Url, handler.fn(r.app))
		case http.MethodDelete:

			r.DELETE(handler.Url, handler.fn(r.app))

		case http.MethodPut:
			r.PUT(handler.Url, handler.fn(r.app))
		default:
			r.Any(handler.Url, handler.fn(r.app))
		}

	}

}

func (r *Controller) initMidleWare() {

	for _, midle := range r.midleWares {
		if ignore := r.IgnoredMidleWare[midle.Name]; !ignore {
			r.Use(midle.fn)
		}
	}
}

func (r *Controller) Start() {

	r.initRouter(r.publicHandlers)
	r.initMidleWare()
	r.initRouter(r.handlers)

	r.Run(fmt.Sprintf(":%d", r.cfg.Port))
}
