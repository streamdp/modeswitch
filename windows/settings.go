package windows

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/streamdp/modeswitch/config"
)

type SettingsWindow struct {
	a fyne.App
	c *config.UserConfig
	fyne.Window
}

func NewSettingsWindow(a fyne.App) *SettingsWindow {
	return &SettingsWindow{
		a: a,
		c: &config.UserConfig{},
		Window: func(a fyne.App) fyne.Window {
			w := a.NewWindow(lang.L("Settings"))
			w.Resize(config.Size)

			return w
		}(a),
	}
}

func (s *SettingsWindow) Create() *SettingsWindow {
	if err := s.c.Load(s.a); err != nil {
		s.showDialogError(err)

		return s
	}

	hostEntry := widget.NewEntry()
	hostEntry.SetPlaceHolder(lang.L("Enter hostname or IP"))
	if s.c.Host != "" {
		hostEntry.SetText(s.c.Host)
	}

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder(lang.L("Enter port number"))
	if s.c.Port != "" {
		portEntry.SetText(s.c.Port)
	}

	isSshCheckBox := widget.NewCheck("", nil)
	isSshCheckBox.Checked = s.c.IsSsh

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder(lang.L("Enter username"))
	if s.c.UserName != "" {
		usernameEntry.SetText(s.c.UserName)
	}

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder(lang.L("Enter password"))
	if s.c.Password != "" {
		passwordEntry.SetText(s.c.Password)
	}

	interfaceId := widget.NewEntry()
	interfaceId.SetPlaceHolder(lang.L("Enter interface Id, such as UsbLte0"))
	if s.c.InterfaceId != "" {
		interfaceId.SetText(s.c.InterfaceId)
	}

	initLteEntry := widget.NewEntry()
	initLteEntry.SetPlaceHolder(lang.L("Enter init string for LTE mode"))
	if s.c.InitLte != "" {
		initLteEntry.SetText(s.c.InitLte)
	}

	initUmtsEntry := widget.NewEntry()
	initUmtsEntry.SetPlaceHolder(lang.L("Enter init string for UMTS mode"))
	if s.c.InitUmts != "" {
		initUmtsEntry.SetText(s.c.InitUmts)
	}

	saveButton := widget.NewButton(lang.L("Save"), func() {
		s.c.UserName = usernameEntry.Text
		s.c.Password = passwordEntry.Text
		s.c.IsSsh = isSshCheckBox.Checked
		s.c.Host = hostEntry.Text
		s.c.Port = portEntry.Text
		s.c.InterfaceId = interfaceId.Text
		s.c.InitLte = initLteEntry.Text
		s.c.InitUmts = initUmtsEntry.Text
		if err := s.c.Save(s.a); err != nil {
			s.showDialogError(err)

			return
		}
		s.showDialogInfo(lang.L("Settings Saved"), lang.L("Settings have been saved"))
	})

	s.SetContent(
		container.NewVScroll(
			container.NewGridWithRows(2,
				container.New(
					layout.NewFormLayout(),
					widget.NewLabel(lang.L("Host")+":"),
					hostEntry,
					widget.NewLabel(lang.L("Port")+":"),
					portEntry,
					widget.NewLabel(lang.L("Ssh")+":"),
					isSshCheckBox,
					widget.NewLabel(lang.L("Username")+":"),
					usernameEntry,
					widget.NewLabel(lang.L("Password")+":"),
					passwordEntry,
					widget.NewLabel(lang.L("Interface")+":"),
					interfaceId,
					widget.NewLabel(lang.L("Init LTE")+":"),
					initLteEntry,
					widget.NewLabel(lang.L("Init UMTS")+":"),
					initUmtsEntry,
				),
				container.NewCenter(
					saveButton,
				),
			),
		),
	)

	return s
}

func (s *SettingsWindow) showDialogError(err error) {
	dialog.ShowError(err, s)
}

func (s *SettingsWindow) showDialogInfo(title, msg string) {
	dialog.ShowInformation(title, msg, s)
}
