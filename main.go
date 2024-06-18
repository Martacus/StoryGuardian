package main

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/application"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	wailsAssetHandler := application.AssetFileServerFS(assets)
	finalAssetHandler := wrapWithImageServer(wailsAssetHandler)

	app := application.New(application.Options{
		Name:        "storyguardian",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(appManager),
			application.NewService(projectManager),
			application.NewService(entityManager),
		},
		Assets: application.AssetOptions{
			Handler: finalAssetHandler,
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

func wrapWithImageServer(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request Path %v \n", r.URL.Path)
		if r.URL.Path == "/" {
			originalHandler.ServeHTTP(w, r)
			return
		}

		if len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/images/" {
			filePath := r.URL.Path[8:]
			fmt.Printf("File Path %v \n", filePath)

			file, err := os.Open(filePath)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {

				}
			}(file)

			contentType := getContentType(filePath)
			w.Header().Set("Content-Type", contentType)

			if _, err := io.Copy(w, file); err != nil {
				http.Error(w, "Failed to read file", http.StatusInternalServerError)
				return
			}

			return
		}

		originalHandler.ServeHTTP(w, r)
	})
}

func getContentType(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream" // fallback content type
	}
}
