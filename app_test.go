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

	service.Register("meizi-handler", reflect.ValueOf(func(id int) string {

		fmt.Println("调用 meizi_handler 成功", id)
		return "ok"
	}))

	var handlers = []*Handler{
		{
			Method: "POST",
			Url:    "/test",
			fn: func(*App) gin.HandlerFunc {
				return func(c *gin.Context) {
					c.JSON(200, map[string]interface{}{"status": "ok"})
				}
			},
			Name: "测试",
		},
		{
			Method: "GET",
			Url:    "/meizi",
			fn: func(app *App) gin.HandlerFunc {
				return func(c *gin.Context) {

					var params struct {
						Id int
					}
					if err := c.ShouldBind(&params); err != nil {
						c.JSON(200, map[string]interface{}{"status": err.Error()})
					}
					result := app.service.Call("meizi-handler", []interface{}{params.Id})
					c.JSON(200, result)

				}
			},
			Name: "妹子",
		},
	}

	controller.WithHandlers(handlers, false)

	app.Start()

}
