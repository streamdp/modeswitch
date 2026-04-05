package windows

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/streamdp/modeswitch/config"
	"github.com/streamdp/modeswitch/device"
)

type SmsWindow struct {
	app fyne.App
	cfg *config.UserConfig
	fyne.Window
}

func NewSmsWindow(app fyne.App) *SmsWindow {
	return &SmsWindow{
		app: app,
		cfg: &config.UserConfig{},
		Window: func(a fyne.App) fyne.Window {
			w := a.NewWindow(lang.L("Sms"))
			w.Resize(config.Size)

			return w
		}(app),
	}
}

func (s *SmsWindow) Create() *SmsWindow {
	outputMemo := widget.NewMultiLineEntry()
	outputMemo.Wrapping = fyne.TextWrapWord
	outputMemo.Disable()

	readSms := widget.NewButton(lang.L("Read"), s.runCmdFn(outputMemo,
		func() {
			if out, err := device.SendCmd(s.cfg, func(d *device.Device) ([]byte, error) {
				return d.SmsList(device.SmsUnread)
			}); err != nil {
				fyne.DoAndWait(func() {
					outputMemo.Append(lang.L("device communication failed") + ": " + err.Error())
				})
				dialog.ShowError(fmt.Errorf("%s:\n%w", lang.L("failed to send commands"), err), s)
			} else {
				fyne.DoAndWait(func() {
					outputMemo.Append(fmt.Sprintf(lang.L("Unread SMS") + ": "))
					outputMemo.Append("\n" + string(out))
				})
			}
		},
	))
	deleteSms := widget.NewButton(lang.L("Delete All"), s.runCmdFn(outputMemo,
		func() {
			if _, err := device.SendCmd(s.cfg, func(d *device.Device) ([]byte, error) {
				_, _ = d.SmsDelete("")

				return nil, nil
			}); err != nil {
				fyne.DoAndWait(func() {
					outputMemo.Append(lang.L("device communication failed") + ": " + err.Error())
				})
				dialog.ShowError(fmt.Errorf("%s:\n%w", lang.L("failed to send commands"), err), s)
			} else {
				fyne.DoAndWait(func() {
					outputMemo.Append(lang.L("SMS has been successfully deleted"))
				})
			}
		},
	))

	s.SetContent(container.NewVSplit(
		container.NewGridWithRows(4,
			layout.NewSpacer(),
			container.NewGridWithColumns(2,
				readSms,
				deleteSms,
			),
			layout.NewSpacer(),
			container.NewCenter(
				widget.NewButton(lang.L("Close"), func() {
					s.Close()
				})),
		),
		container.NewStack(outputMemo),
	))

	return s
}

func (s *SmsWindow) runCmdFn(output *widget.Entry, f func()) func() {
	return func() {
		if err := loadConfig(s.app, s.cfg, output); err != nil {
			dialog.ShowError(err, s)

			return
		}

		output.Append(lang.L("trying to connect") + "...\n")
		go func() {
			f()
		}()
	}
}
