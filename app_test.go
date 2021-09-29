package app

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestApp(t *testing.T) {

	var cfg = &Config{
		ControllerCfg: &ControllerConfig{
			Port: 8081,
		},
		// DaoCfg: &DaoConfig{
		// 	UserName:     "user",
		// 	Password:     "123",
		// 	Host:         "127.0.0.1",
		// 	Port:         3306,
		// 	DatabaseName: "test",
		// },
	}
	var app = NewApp(cfg)

	// logic handlers
	app.Register("test-handler", reflect.ValueOf(func(id int) string {

		fmt.Println("调用 test_handler 成功", id)

		return "ok"
	}))

	// dispatcher
	var handlers = []*Handler{
		{
			Method: "GET",
			Url:    "/test",
			Name:   "测试",
			Fn: func(app *App) gin.HandlerFunc {
				return func(c *gin.Context) {

					var params struct {
						Id int
					}
					if err := c.ShouldBind(&params); err != nil {
						c.JSON(200, map[string]interface{}{"status": err.Error()})
					}
					result := app.Call("test-handler", []interface{}{params.Id})
					c.JSON(200, result)

				}
			},
		},
	}

	app.Controller.WithHandlers(handlers, false)

	app.Start()

}
