package windows

import (
	"errors"
	"fmt"

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

type SmsWindow struct {
	a fyne.App
	c *config.UserConfig
	fyne.Window
}

func NewSmsWindow(a fyne.App) *SmsWindow {
	return &SmsWindow{
		a: a,
		c: &config.UserConfig{},
		Window: func(a fyne.App) fyne.Window {
			w := a.NewWindow(lang.L("Sms"))
			w.Resize(config.Size)

			return w
		}(a),
	}
}

func (s *SmsWindow) Create() *SmsWindow {
	outputMemo := widget.NewMultiLineEntry()
	outputMemo.Wrapping = fyne.TextWrapWord
	outputMemo.Disable()

	readSms := widget.NewButton(lang.L("Read"), s.runCmdFn(outputMemo,
		func() {
			if out, err := s.sendSMSListCmd(device.SmsUnread); err != nil {
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
			if err := s.sendDeleteSMSCmd(); err != nil {
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
		output.SetText(lang.L("load config") + "...\n")

		if err := s.c.Load(s.a); err != nil {
			output.Append(lang.L("load config error") + ": " + err.Error())
			dialog.ShowError(err, s)

			return
		}

		if s.c.Host == "" || s.c.Port == "" || s.c.InterfaceId == "" {
			errConfig := errors.New(lang.L("you need to configure the app first"))
			output.Append(lang.L("config error") + ": " + errConfig.Error())
			dialog.ShowError(errConfig, s)

			return
		}

		output.Append(lang.L("trying to connect") + "...\n")
		go func() {
			f()
		}()
	}
}

func (s *SmsWindow) sendSMSListCmd(args ...string) ([]byte, error) {
	var (
		sender device.Sender
		err    error
	)

	if s.c.IsSsh {
		if sender, err = ssh_client.New(s.c); err != nil {
			return nil, err
		}
	} else {
		if sender, err = telnet_client.New(s.c); err != nil {
			return nil, err
		}
	}
	defer func(sender device.Sender) {
		_ = sender.Close()
	}(sender)

	return device.SmsList(sender, s.c, args...)
}

func (s *SmsWindow) sendDeleteSMSCmd() error {
	return fmt.Errorf("unimplemented")
}
