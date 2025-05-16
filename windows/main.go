package windows

import (
	"errors"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/streamdp/modeswitch/config"
	"github.com/streamdp/modeswitch/connections"
)

type MainWindow struct {
	a fyne.App
	c *config.UserConfig
	fyne.Window
}

func NewMainWindow(a fyne.App) *MainWindow {
	return &MainWindow{
		a: a,
		c: &config.UserConfig{},
		Window: func(a fyne.App) fyne.Window {
			w := a.NewWindow(lang.L("Mode switch"))
			w.Resize(config.Size)

			return w
		}(a),
	}
}

func (m *MainWindow) Create() *MainWindow {
	outputMemo := widget.NewMultiLineEntry()
	outputMemo.Wrapping = fyne.TextWrapWord
	outputMemo.Disable()

	switchToLteButton := widget.NewButton(lang.L("Switch to LTE"), m.switchButtonFn(config.LTE, outputMemo))
	switchToUMTSButton := widget.NewButton(lang.L("Switch to UMTS"), m.switchButtonFn(config.UMTS, outputMemo))

	m.SetContent(container.NewVSplit(
		container.NewGridWithRows(4,
			layout.NewSpacer(),
			container.NewGridWithColumns(2,
				switchToLteButton,
				switchToUMTSButton,
			),
			layout.NewSpacer(),
			container.NewCenter(
				widget.NewButton(lang.L("Quit"), func() {
					m.Close()
					os.Exit(0)
				})),
		),
		container.NewStack(outputMemo),
	))

	return m
}

func (m *MainWindow) switchButtonFn(mode string, output *widget.Entry) func() {
	return func() {
		output.SetText(lang.L("load config") + "...\n")

		if err := m.c.Load(m.a); err != nil {
			output.Append(lang.L("load config error") + ": " + err.Error())
			m.showDialogError(err)

			return
		}

		if m.c.Host == "" || m.c.Port == "" || m.c.InterfaceId == "" {
			errConfig := errors.New(lang.L("you need to configure the app first"))
			output.Append(lang.L("config error") + ": " + errConfig.Error())
			m.showDialogError(errConfig)

			return
		}

		output.Append(lang.L("trying to connect") + "...\n")
		go fyne.Do(func() {
			if err := m.sendCommand(mode); err != nil {
				output.Append(lang.L("connection error") + ": " + err.Error())
				m.showDialogError(fmt.Errorf("%s:\n%w", lang.L("unable to connect"), err))
			} else {
				output.Append(fmt.Sprintf(lang.L("commands to switch %s mode were sent successfully", mode)))
				output.Append("\n" + lang.L("wait few minutes, and check internet connection"))
			}
		})
	}
}

func (m *MainWindow) showDialogError(err error) {
	dialog.ShowError(err, m)
}

func (m *MainWindow) sendCommand(mode string) error {
	if m.c.IsSsh {
		return connections.SendSshCommand(mode, m.c)
	}

	return connections.SendTelnetCommand(mode, m.c)
}
