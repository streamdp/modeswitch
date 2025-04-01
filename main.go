package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/streamdp/modeswitch/windows"
)

func main() {
	a := app.New()

	m := windows.NewMainWindow(a).Create()
	m.SetMainMenu(windows.CreateMainMenu(a))
	m.ShowAndRun()
}
