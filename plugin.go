package app

type Plugin interface {
	Initialize(app *App)
}

