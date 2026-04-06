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
			w := a.NewWindow(lang.L("SMS"))
			w.Resize(config.Size)

			return w
		}(app),
	}
}

func (s *SmsWindow) Create() *SmsWindow {
	outputMemo := widget.NewMultiLineEntry()
	outputMemo.Wrapping = fyne.TextWrapWord

	readSms := widget.NewButton(lang.L("Read"), s.runCmdFn(outputMemo,
		func() {
			resp, err := device.SendCmd(s.cfg, func(d *device.Device) (*device.Response, error) {
				return d.SmsList(device.SmsUnread)
			})
			if err != nil {
				fyne.DoAndWait(func() {
					outputMemo.Append(lang.L("device communication failed") + ": " + err.Error())
				})
				dialog.ShowError(fmt.Errorf("%s:\n%w", lang.L("failed to send commands"), err), s)

				return
			}

			fyne.DoAndWait(func() {
				outputMemo.Append(fmt.Sprintf(lang.L("Unread SMS") + ": "))

				for _, msg := range resp.SMS.Messages {
					outputMemo.Append(fmt.Sprintf("\n%v\n", msg.String()))
				}
			})
		},
	))
	deleteSms := widget.NewButton(lang.L("Delete All"), s.runCmdFn(outputMemo,
		func() {
			_, err := device.SendCmd(s.cfg, func(d *device.Device) (*device.Response, error) {
				resp, err := d.SmsList(device.SmsNoContent)
				if err != nil {
					return nil, err
				}

				for _, msg := range resp.SMS.Messages {
					resp, err = d.SmsDelete(msg.ID)
					if err != nil {
						return nil, err
					}
					if resp.Text != "" {
						fyne.DoAndWait(func() { outputMemo.Append(resp.Text) })
					}
				}

				return nil, nil
			})
			if err != nil {
				fyne.DoAndWait(func() {
					outputMemo.Append(lang.L("device communication failed") + ": " + err.Error())
				})
				dialog.ShowError(fmt.Errorf("%s:\n%w", lang.L("failed to send commands"), err), s)

				return
			}

			fyne.DoAndWait(func() {
				outputMemo.Append("\n" + lang.L("SMS has been successfully deleted"))
			})
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
