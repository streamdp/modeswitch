package windows

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/streamdp/modeswitch/config"
)

func CreateSettingsWindow(a fyne.App) (w fyne.Window) {
	w = a.NewWindow("Settings")

	w.Resize(config.Size)

	c := &config.UserConfig{}
	if err := c.Load(a); err != nil {
		dialog.ShowError(err, w)
	}

	hostEntry := widget.NewEntry()
	hostEntry.SetPlaceHolder("Enter hostname or IP")
	if c.Host != "" {
		hostEntry.SetText(c.Host)
	}

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Enter port number")
	if c.Port != "" {
		portEntry.SetText(c.Port)
	}

	isSshCheckBox := widget.NewCheck("", nil)
	isSshCheckBox.Checked = c.IsSsh

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Enter username")
	if c.UserName != "" {
		usernameEntry.SetText(c.UserName)
	}

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter password")
	if c.Password != "" {
		passwordEntry.SetText(c.Password)
	}

	interfaceId := widget.NewEntry()
	interfaceId.SetPlaceHolder("Enter interface Id, such as UsbLte0")
	if c.InterfaceId != "" {
		interfaceId.SetText(c.InterfaceId)
	}

	initLteEntry := widget.NewEntry()
	initLteEntry.SetPlaceHolder("Enter init string for LTE mode")
	if c.InitLte != "" {
		initLteEntry.SetText(c.InitLte)
	}

	initUmtsEntry := widget.NewEntry()
	initUmtsEntry.SetPlaceHolder("Enter init string for UMTS mode")
	if c.InitUmts != "" {
		initUmtsEntry.SetText(c.InitUmts)
	}

	saveButton := widget.NewButton("Save", func() {
		c.UserName = usernameEntry.Text
		c.Password = passwordEntry.Text
		c.IsSsh = isSshCheckBox.Checked
		c.Host = hostEntry.Text
		c.Port = portEntry.Text
		c.InterfaceId = interfaceId.Text
		c.InitLte = initLteEntry.Text
		c.InitUmts = initUmtsEntry.Text
		if err := c.Save(a); err != nil {
			dialog.ShowError(err, w)
		}
		dialog.ShowInformation("Settings Saved", "Connection preferences and credentials have been saved.", w)
		w.Close()
	})

	w.SetContent(
		container.NewVBox(
			container.New(
				layout.NewFormLayout(),
				widget.NewLabel("Host:"),
				hostEntry,
				widget.NewLabel("Port:"),
				portEntry,
				widget.NewLabel("Ssh:"),
				isSshCheckBox,
				widget.NewLabel("Username:"),
				usernameEntry,
				widget.NewLabel("Password:"),
				passwordEntry,
				widget.NewLabel("Interface"),
				interfaceId,
				widget.NewLabel("init LTE:"),
				initLteEntry,
				widget.NewLabel("init UMTS:"),
				initUmtsEntry,
			),
			layout.NewSpacer(),
			saveButton,
		))

	return
}
