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
)

var errConfig = errors.New(lang.L("you need to configure the app first"))

type MainWindow struct {
	app fyne.App
	cfg *config.UserConfig
	fyne.Window
}

func NewMainWindow(app fyne.App) *MainWindow {
	return &MainWindow{
		app: app,
		cfg: &config.UserConfig{},
		Window: func(a fyne.App) fyne.Window {
			w := a.NewWindow(lang.L("Mode switch"))
			w.Resize(config.Size)

			return w
		}(app),
	}
}

func (m *MainWindow) Create() *MainWindow {
	outputMemo := widget.NewMultiLineEntry()
	outputMemo.Wrapping = fyne.TextWrapWord
	// outputMemo.Disable()

	switchToLteButton := widget.NewButton(lang.L("Switch to LTE"), m.switchModemModeFn(config.LTE, outputMemo))
	switchToUMTSButton := widget.NewButton(lang.L("Switch to UMTS"), m.switchModemModeFn(config.UMTS, outputMemo))

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

func (m *MainWindow) switchModemModeFn(mode string, output *widget.Entry) func() {
	return func() {
		if err := loadConfig(m.app, m.cfg, output); err != nil {
			dialog.ShowError(err, m)

			return
		}

		output.Append(lang.L("trying to connect") + "...\n")
		go func() {
			_, err := device.SendCmd(m.cfg, func(d *device.Device) (*device.Response, error) {
				return d.SwitchLTEModemMode(mode)
			})
			if err != nil {
				fyne.DoAndWait(func() { output.Append(lang.L("device communication failed") + ": " + err.Error()) })
				dialog.ShowError(fmt.Errorf("%s:\n%w", lang.L("failed to send commands"), err), m)

				return
			}

			fyne.DoAndWait(func() {
				output.Append(fmt.Sprintf(lang.L("commands to switch %s mode were sent successfully"), mode))
				output.Append("\n" + lang.L("wait few minutes, and check internet connection"))
			})
		}()
	}
}

func loadConfig(app fyne.App, cfg *config.UserConfig, output *widget.Entry) error {
	output.SetText(lang.L("load config") + "...\n")

	if err := cfg.Load(app); err != nil {
		output.Append(lang.L("load config error") + ": " + err.Error())

		return err
	}

	if cfg.Host == "" || cfg.Port == "" || cfg.InterfaceId == "" {
		output.Append(lang.L("config error") + ": " + errConfig.Error())

		return errConfig
	}

	return nil
}
