package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/streamdp/modeswitch/windows"
)

func main() {
	a := app.New()
	w := windows.CreateMainWindow(a)
	w.SetMainMenu(windows.CreateMainMenu(a))
	w.ShowAndRun()
}
