package config

import (
	"time"

	"fyne.io/fyne/v2"
	"github.com/streamdp/modeswitch/encryption"
)

const (
	UMTS = "umts"
	LTE  = "lte"

	DefaultTimeout = 5 * time.Second
)

var Size = fyne.Size{
	Width:  240,
	Height: 480,
}

type UserConfig struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	IsSsh       bool   `json:"is_ssh"`
	InterfaceId string `json:"interface_id"`
	InitLte     string `json:"init_lte"`
	InitUmts    string `json:"init_umts"`
}

func (uc *UserConfig) Save(a fyne.App) (err error) {
	a.Preferences().SetString("host", uc.Host)
	a.Preferences().SetString("port", uc.Port)
	a.Preferences().SetBool("is_ssh", uc.IsSsh)
	a.Preferences().SetString("username", uc.UserName)
	var password string
	if password, err = encryption.Encrypt(uc.Password, uc.UserName); err != nil {
		return
	}
	a.Preferences().SetString("password", password)
	a.Preferences().SetString("init_lte", uc.InitLte)
	a.Preferences().SetString("init_umts", uc.InitUmts)
	a.Preferences().SetString("interface_id", uc.InterfaceId)
	return
}

func (uc *UserConfig) Load(a fyne.App) (err error) {
	uc.UserName = a.Preferences().String("username")
	var password = a.Preferences().String("password")
	if password, err = encryption.Decrypt(password, uc.UserName); err != nil {
		return
	}
	uc.Password = password
	uc.Host = a.Preferences().String("host")
	uc.Port = a.Preferences().String("port")
	uc.IsSsh = a.Preferences().Bool("is_ssh")
	uc.InitLte = a.Preferences().String("init_lte")
	uc.InitUmts = a.Preferences().String("init_umts")
	uc.InterfaceId = a.Preferences().String("interface_id")
	return nil
}
