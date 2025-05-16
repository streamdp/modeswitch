package main

import (
	"embed"
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/lang"
	"github.com/streamdp/modeswitch/windows"
)

//go:embed translations
var translations embed.FS

func main() {
	if err := lang.AddTranslationsFS(translations, "translations"); err != nil {
		log.Fatal("failed to load translations")
	}

	a := app.New()

	m := windows.NewMainWindow(a).Create()
	m.SetMainMenu(windows.CreateMainMenu(a))
	m.ShowAndRun()
}
