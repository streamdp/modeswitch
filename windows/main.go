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
	"github.com/streamdp/modeswitch/device"
	"github.com/streamdp/modeswitch/pkg/ssh_client"
	"github.com/streamdp/modeswitch/pkg/telnet_client"
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
			dialog.ShowError(err, m)

			return
		}

		if m.c.Host == "" || m.c.Port == "" || m.c.InterfaceId == "" {
			errConfig := errors.New(lang.L("you need to configure the app first"))
			output.Append(lang.L("config error") + ": " + errConfig.Error())
			dialog.ShowError(errConfig, m)

			return
		}

		output.Append(lang.L("trying to connect") + "...\n")
		go func() {
			if err := m.sendSwitchModeCmd(mode); err != nil {
				fyne.DoAndWait(func() { output.Append(lang.L("connection error") + ": " + err.Error()) })
				dialog.ShowError(fmt.Errorf("%s:\n%w", lang.L("unable to connect"), err), m)
			} else {
				fyne.DoAndWait(func() {
					output.Append(fmt.Sprintf(lang.L("commands to switch %s mode were sent successfully"), mode))
					output.Append("\n" + lang.L("wait few minutes, and check internet connection"))
				})
			}
		}()
	}
}

func (m *MainWindow) sendSwitchModeCmd(mode string) error {
	var (
		sender device.Sender
		err    error
	)

	if m.c.IsSsh {
		if sender, err = ssh_client.New(m.c); err != nil {
			return err
		}
	} else {
		if sender, err = telnet_client.New(m.c); err != nil {
			return err
		}
	}
	defer func(sender device.Sender) {
		_ = sender.Close()
	}(sender)

	return device.SwitchLTEModemMode(sender, mode, m.c)
}
