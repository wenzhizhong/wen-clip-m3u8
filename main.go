package main

import (
	"clipM3u8Media/goApi"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	m3u8HandlerApi := &goApi.M3u8Handler{
		Ctx: &app.ctx,
	}
	runtimeApi := &goApi.Runtime{
		Ctx: &app.ctx,
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "clipM3u8Media",
		Width:         1024,
		Height:        768,
		OnBeforeClose: (&goApi.Runtime{}).BeforeClose,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			m3u8HandlerApi,
			runtimeApi,
		},
		Windows: &windows.Options{
			DisablePinchZoom:     true,  // 禁用所有缩放控制（包括快捷键和触控板手势）
			IsZoomControlEnabled: false, // 禁止用户改变缩放因子（关键选项）
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
