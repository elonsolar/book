package app

import "fmt"

type App struct {
	controller *Controller
	service    *Service
	dao        *Dao
	errors     []error
}

func NewApp() *App {

	return &App{}
}

func (a *App) Start() {

	a.controller.Start()
}

func (a *App) Error() {

	fmt.Println(a.errors)
}
