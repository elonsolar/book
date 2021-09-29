package controller

import (
	. "github.com/elonsolar/easy-app"
	"github.com/gin-gonic/gin"
)

// router
func AddUser(app *App) gin.HandlerFunc {

	return func(c *gin.Context) {
		var params struct {
			Id int
		}
		if err := c.ShouldBind(&params); err != nil {
			c.JSON(200, map[string]interface{}{"status": err.Error()})
		}
		result := app.Call("logic_add_user", []interface{}{params.Id})
		c.JSON(200, result)
	}
}
