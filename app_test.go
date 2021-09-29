package app

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestApp(t *testing.T) {

	var app = NewApp()
	var controller = NewController(&ControllerConfig{Port: 8081}, app)
	var service = NewService(app)
	app.controller = controller
	app.service = service

	app.Register("test-handler", reflect.ValueOf(func(id int) string {

		fmt.Println("调用 test_handler 成功", id)
		return "ok"
	}))

	var handlers = []*Handler{
		{
			Method: "GET",
			Url:    "/test",
			Name:   "测试",
			fn: func(app *App) gin.HandlerFunc {
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

	controller.WithHandlers(handlers, false)

	app.Start()

}
