package service

import (
	"reflect"

	. "github.com/elonsolar/easy-app"
)

func InitLogicHandler(app *App) {
	app.Register("logic_add_user", reflect.ValueOf(AddUser))

}
