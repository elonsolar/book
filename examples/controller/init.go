package controller

import . "github.com/elonsolar/easy-app"

func InitDispatcherHandler(app *App) {

	app.Controller.AddHandler("测试", "GET", "/test", AddUser, true)
	app.Controller.AddHandler("测试2", "GET", "/test/v2", AddUser, true)

}
