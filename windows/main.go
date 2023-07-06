package windows

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/streamdp/modeswitch/config"
	"github.com/streamdp/modeswitch/connections"
)

func CreateMainWindow(a fyne.App) (w fyne.Window) {
	w = a.NewWindow("Mode switch")

	w.Resize(config.Size)

	outputMemo := widget.NewMultiLineEntry()

	switchToLteButton := widget.NewButton("Switch to LTE", func() {
		go func() {
			if msg, err := switchTo(config.LTE, a); err != nil {
				dialog.ShowError(err, w)
			} else {
				outputMemo.SetText(msg)
			}
		}()
	})

	switchToUMTSButton := widget.NewButton("Switch to UMTS", func() {
		go func() {
			if msg, err := switchTo(config.UMTS, a); err != nil {
				dialog.ShowError(err, w)
			} else {
				outputMemo.SetText(msg)
			}
		}()
	})

	w.SetContent(container.NewVBox(
		layout.NewSpacer(),
		switchToLteButton,
		switchToUMTSButton,
		layout.NewSpacer(),
		widget.NewButton("Quit", func() {
			w.Close()
			os.Exit(0)
		}),
		layout.NewSpacer(),
		outputMemo,
	))

	return
}

func switchTo(mode string, a fyne.App) (_ string, err error) {
	var c = &config.UserConfig{}
	if err = c.Load(a); err != nil {
		return
	}
	if c.Host == "" || c.Port == "" || c.InterfaceId == "" {
		return "", fmt.Errorf("you need to configure the app first.\n specify host, port, interface id at least")
	}
	if err = sendCommand(mode, c); err != nil {
		return
	}
	return fmt.Sprintf("commands to switch %s mode were sent successfully", mode), nil
}

func sendCommand(mode string, c *config.UserConfig) (err error) {
	if c.IsSsh {
		return connections.SendSshCommand(mode, c)
	}
	return connections.SendTelnetCommand(mode, c)
}
