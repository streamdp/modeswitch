package windows

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
	"github.com/streamdp/modeswitch/config"
)

func CreateMainMenu(a fyne.App) *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu(lang.L("Menu"),
			fyne.NewMenuItem(lang.L("Settings"), func() {
				sw := NewSettingsWindow(a).Create()
				sw.Resize(config.Size)
				sw.CenterOnScreen()
				sw.Show()
			}),
		),
	)
}
