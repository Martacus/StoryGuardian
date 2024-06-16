package main

import (
	"embed"
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"log"
	"storyguardian/project"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed frontend/dist
var assets embed.FS

func main() {
	appManager := project.NewApplicationManager()
	projectManager := project.NewProjectManager(appManager)
	entityManager := project.NewEntityManager(projectManager)

	app := application.New(application.Options{
		Name:        "storyguardian",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(appManager),
			application.NewService(projectManager),
			application.NewService(entityManager),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		Height:           720,
		Width:            1280,
	})

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	//go func() {
	//	for {
	//		now := time.Now().Format(time.RFC1123)
	//		app.Events.Emit(&application.WailsEvent{
	//			Name: "time",
	//			Data: now,
	//		})
	//		time.Sleep(time.Second)
	//	}
	//}()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
