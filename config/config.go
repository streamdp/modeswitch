package config

import (
	"fyne.io/fyne/v2"
	"github.com/streamdp/modeswitch/encryption"
)

const (
	UMTS = "umts"
	LTE  = "lte"
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

func (s *UserConfig) Save(a fyne.App) (err error) {
	a.Preferences().SetString("host", s.Host)
	a.Preferences().SetString("port", s.Port)
	a.Preferences().SetBool("is_ssh", s.IsSsh)
	a.Preferences().SetString("username", s.UserName)
	var password string
	if password, err = encryption.Encrypt(s.Password, s.UserName); err != nil {
		return
	}
	a.Preferences().SetString("password", password)
	a.Preferences().SetString("init_lte", s.InitLte)
	a.Preferences().SetString("init_umts", s.InitUmts)
	a.Preferences().SetString("interface_id", s.InterfaceId)
	return
}

func (s *UserConfig) Load(a fyne.App) (err error) {
	s.UserName = a.Preferences().String("username")
	var password = a.Preferences().String("password")
	if password, err = encryption.Decrypt(password, s.UserName); err != nil {
		return
	}
	s.Password = password
	s.Host = a.Preferences().String("host")
	s.Port = a.Preferences().String("port")
	s.IsSsh = a.Preferences().Bool("is_ssh")
	s.InitLte = a.Preferences().String("init_lte")
	s.InitUmts = a.Preferences().String("init_umts")
	s.InterfaceId = a.Preferences().String("interface_id")
	return nil
}
