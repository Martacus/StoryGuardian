package main

import (
	"embed"
	"log"

	"litguardian/internal"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "litguardian",
		Description: "Local-first worldbuilding app for LitRPG and fantasy authors",
		Services: []application.Service{
			application.NewService(&internal.AppConfigService{}),
			application.NewService(&internal.WorldService{}),
			application.NewService(&internal.LayoutService{}),
			application.NewService(&internal.ImageService{}),
			application.NewService(&internal.EntityService{}),
			application.NewService(&internal.LinkService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "LitGuardian",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
