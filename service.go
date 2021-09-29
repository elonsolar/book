package app

// Service mgr all logical handlers
type Service struct {
	app *App
}

func NewService(app *App) *Service {
	return &Service{
		app: app,
	}
}
