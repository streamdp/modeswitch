package main

import (
	"embed"
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"github.com/streamdp/modeswitch/internal/pkg/custom_theme"
	"github.com/streamdp/modeswitch/internal/windows"
)

//go:embed internal/translations
var translations embed.FS

func main() {
	if err := lang.AddTranslationsFS(translations, "internal/translations"); err != nil {
		log.Fatal("failed to load translations")
	}

	a := app.New()
	a.Settings().SetTheme(&custom_theme.Theme{
		Theme: theme.DefaultTheme(),
	})

	m := windows.NewMainWindow(a).Create()
	m.SetMainMenu(windows.CreateMainMenu(a))
	m.ShowAndRun()
}
