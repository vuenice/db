//go:build desktop

package main

import (
	"log"
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func runDesktopApp(handler http.Handler) {
	wailsApp := application.New(application.Options{
		Name:        "chatdb",
		Description: "chatdb Wails application",
		Assets: application.AssetOptions{
			Handler: handler,
		},
	})

	go func() {
		err := wailsApp.Run()
		if err != nil {
			log.Printf("wails app run: %v", err)
		}
	}()
}
