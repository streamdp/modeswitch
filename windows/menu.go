package windows

import (
	"fyne.io/fyne/v2"
	"github.com/streamdp/modeswitch/config"
)

func CreateMainMenu(a fyne.App) *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("Menu",
			fyne.NewMenuItem("Settings", func() {
				sw := NewSettingsWindow(a).Create()
				sw.Resize(config.Size)
				sw.CenterOnScreen()
				sw.Show()
			}),
		),
	)
}
