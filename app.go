package main

import (
	"clipM3u8Media/backend"
	"context"
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// App struct
type App struct {
	app *application.App
}

// NewApp creates a new App application struct
func NewApp(app *application.App) *App {
	return &App{app: app}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Web服务端
	go func() {
		backend.Backend()
	}()
	return nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
