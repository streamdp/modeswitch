package config

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/streamdp/modeswitch/encryption"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestUserConfig_Save(t *testing.T) {
	a := app.NewWithID("com.test.testId")

	testCases := []UserConfig{
		{
			UserName:    "username",
			Password:    "password",
			Host:        "192.168.100.1",
			Port:        "22",
			IsSsh:       true,
			InterfaceId: "UsbLte1",
			InitLte:     "at+xact=2,,,107",
			InitUmts:    "at+xact=1,,,0",
		},
		{
			UserName:    "username2",
			Password:    "password2",
			Host:        "192.168.1.1",
			Port:        "23",
			IsSsh:       false,
			InterfaceId: "UsbLte0",
			InitLte:     "at+xact=2,,,103,107",
			InitUmts:    "at+xact=1,,,101,102",
		},
	}

	for n, tt := range testCases {
		tt := tt
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			err := tt.Save(a)
			assert.NoError(t, err)

			decryptedPassword, err := encryption.Decrypt(a.Preferences().String("password"), tt.UserName)
			assert.NoError(t, err)

			assert.Equal(t, tt.UserName, a.Preferences().String("username"))
			assert.Equal(t, tt.Password, decryptedPassword)
			assert.Equal(t, tt.Host, a.Preferences().String("host"))
			assert.Equal(t, tt.Port, a.Preferences().String("port"))
			assert.Equal(t, tt.IsSsh, a.Preferences().Bool("is_ssh"))
			assert.Equal(t, tt.InterfaceId, a.Preferences().String("interface_id"))
			assert.Equal(t, tt.InitLte, a.Preferences().String("init_lte"))
			assert.Equal(t, tt.InitUmts, a.Preferences().String("init_umts"))
		})
	}
}

func TestUserConfig_Load(t *testing.T) {
	a := app.NewWithID("com.test.testId")

	testCases := []UserConfig{
		{
			UserName:    "betterUsername",
			Password:    "betterPassword",
			Host:        "10.4.10.11",
			Port:        "22",
			IsSsh:       true,
			InterfaceId: "UsbModem0",
			InitLte:     "at+mtsm=1",
			InitUmts:    "at+gtpkgver?",
		},
		{
			UserName:    "nextUsername",
			Password:    "securePassword",
			Host:        "3.54.12.222",
			Port:        "23",
			IsSsh:       false,
			InterfaceId: "UsbModem0",
			InitLte:     "at@nvm:cal_usbmode.num=0",
			InitUmts:    "at@nvm:cal_usbmode.num=10",
		},
	}

	for n, tt := range testCases {
		tt := tt
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			assert.NoError(t, tt.Save(a))

			c := &UserConfig{}
			assert.NoError(t, c.Load(a))

			assert.Equal(t, c.UserName, tt.UserName)
			assert.Equal(t, c.Password, tt.Password)
			assert.Equal(t, c.Host, tt.Host)
			assert.Equal(t, c.Port, tt.Port)
			assert.Equal(t, c.IsSsh, tt.IsSsh)
			assert.Equal(t, c.InterfaceId, tt.InterfaceId)
			assert.Equal(t, c.InitLte, tt.InitLte)
			assert.Equal(t, c.InitUmts, tt.InitUmts)
		})
	}
}

func TestConstants(t *testing.T) {
	assert.Equal(t, UMTS, "umts")
	assert.Equal(t, LTE, "lte")
	assert.Equal(t, Size, fyne.Size{
		Width:  240,
		Height: 480,
	})
}
