package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:            "CP77 Mod Conflict Manager",
		Width:            1100,
		Height:           650,
		MinWidth:         900,
		MinHeight:        550,
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup:        app.startup,
		Bind:             []interface{}{app},
		BackgroundColour: &options.RGBA{R: 13, G: 13, B: 13, A: 255},
	})
	if err != nil {
		panic(err)
	}
}
