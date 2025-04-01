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

var errConfig = fmt.Errorf("you need to configure the app first.\nspecify host, port, interface id at least")

type MainWindow struct {
	a fyne.App
	c *config.UserConfig
	fyne.Window
}

func NewMainWindow(a fyne.App) *MainWindow {
	return &MainWindow{
		a: a,
		c: &config.UserConfig{},
		Window: func(a fyne.App) (w fyne.Window) {
			w = a.NewWindow("Mode switch")
			w.Resize(config.Size)
			return
		}(a),
	}
}

func (m *MainWindow) Create() fyne.Window {
	outputMemo := widget.NewMultiLineEntry()
	outputMemo.Wrapping = fyne.TextWrapWord

	switchToLteButton := widget.NewButton("Switch to LTE", m.switchButtonFn(config.LTE, outputMemo))
	switchToUMTSButton := widget.NewButton("Switch to UMTS", m.switchButtonFn(config.UMTS, outputMemo))

	m.SetContent(container.NewVBox(
		layout.NewSpacer(),
		switchToLteButton,
		switchToUMTSButton,
		layout.NewSpacer(),
		widget.NewButton("Quit", func() {
			m.Close()
			os.Exit(0)
		}),
		layout.NewSpacer(),
		outputMemo,
	))

	return m
}

func (m *MainWindow) switchButtonFn(mode string, output *widget.Entry) func() {
	return func() {
		output.SetText("load config...\n")

		if err := m.c.Load(m.a); err != nil {
			output.Append(fmt.Sprintf("load config error: %s\n", err))
			m.showDialogError(err)

			return
		}

		if m.c.Host == "" || m.c.Port == "" || m.c.InterfaceId == "" {
			output.Append(fmt.Sprintf("config error: %v\n", errConfig))
			m.showDialogError(errConfig)

			return
		}

		go func() {
			output.Append("trying to connect...\n")

			if err := m.sendCommand(mode); err != nil {
				output.Append(fmt.Sprintf("connection error: %v\n", err))
				m.showDialogError(err)
			} else {
				output.Append(fmt.Sprintf("commands to switch %s mode were sent successfully.\nwait few minutes,"+
					" and check internet connection.", mode),
				)
			}
		}()
	}
}

func (m *MainWindow) showDialogError(err error) {
	dialog.ShowError(err, m)
}

func (m *MainWindow) sendCommand(mode string) (err error) {
	if m.c.IsSsh {
		return connections.SendSshCommand(mode, m.c)
	}
	return connections.SendTelnetCommand(mode, m.c)
}
