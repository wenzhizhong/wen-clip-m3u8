package main

import (
	"clipM3u8Media/goApi"
	"clipM3u8Media/goApi/common"
	"context"
	"embed"
	"fmt"

	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS
var app *application.App

func init() {
	// Register a custom event whose associated data type is string.
	// This is not required, but the binding generator will pick up registered events
	// and provide a strongly typed JS/TS API for them.
	application.RegisterEvent[string]("time")
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	ctx := context.Background()
	m3u8HandlerApi := &goApi.M3u8Handler{
		Ctx: &ctx,
	}
	getConstantApi := &goApi.GetConstant{
		Ctx: &ctx,
	}
	// runtimeApi := &goApi.Runtime{
	// 	App: app,
	// }

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app = application.New(application.Options{
		Name:        "clipM3u8Media",
		Description: "本地m3u8流媒体剪辑。Local m3u8 streaming file media editing.",
		Services:    []application.Service{
			// application.NewService(&GreetService{}),
			// application.NewService(m3u8HandlerApi),
			// application.NewService(runtimeApi),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		// Windows: application.WindowsOptions{
		// 	DisableQuitOnLastWindowClosed: false,
		// },
		// Linux: application.LinuxOptions{
		// 	DisableQuitOnLastWindowClosed: false,
		// },
		// ShouldQuit: func() bool {
		// 	return false
		// },
	})

	app.RegisterService(application.NewService(NewApp(app)))
	app.RegisterService(application.NewService(m3u8HandlerApi))
	app.RegisterService(application.NewService(getConstantApi))
	app.RegisterService(application.NewService(&goApi.Runtime{App: app}))

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "clipM3u8Media",
		Name:   common.AppMainWindowName,
		URL:    "/",
		Width:  1024,
		Height: 768,
		// OnBeforeClose: (&goApi.Runtime{}).BeforeClose,
		// OnStartup:        app.startup,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
	})

	preloadVideoWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "预加载",
		Name:   common.AppPreUploadWindowName,
		URL:    "/#/preloadVideo",
		Width:  1200,
		Height: 400,
		// OnBeforeClose: (&goApi.Runtime{}).BeforeClose,
		// OnStartup:        app.startup,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		Hidden:           true,
		InitialPosition:  1,
		X:                300,
		Y:                300,
	})

	mainWindow.RegisterHook(events.Common.WindowClosing, exitMainWindowConfirm)
	preloadVideoWindow.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) { toggleWindow(event, common.AppPreUploadWindowName) })
	preloadVideoWindow.Hide()

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}

func exitMainWindowConfirm(event *application.WindowEvent) {
	// event.Cancel()
	defBtn := &application.Button{
		Label: "Yes",
		// Label:     "OK",
		// Label:     "是",
		IsDefault: true,
		Callback: func() {
			app.Quit()
		},
	}
	cancelBtn := &application.Button{
		Label: "No",
		// Label:    "Cancel",
		// Label:    "否",
		IsCancel: true,
		Callback: func() {
			event.Cancel()
		},
	}

	buttons := []*application.Button{}
	buttons = append(buttons, defBtn)
	buttons = append(buttons, cancelBtn)

	app.Dialog.Question().
		SetTitle("提示").
		SetMessage("是否退出?").
		AddButtons(buttons).
		Show()
}

func toggleWindow(event *application.WindowEvent, windowName string) {
	event.Cancel()
	window, ok := app.Window.GetByName(windowName)
	if !ok {
		fmt.Printf("toggleWindow(): window %s not found", windowName)
		return
	}
	if window.IsVisible() {
		window.Hide()
		return
	}
	window.Show()
}
