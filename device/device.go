package device

import (
	"github.com/streamdp/modeswitch/config"
	"github.com/streamdp/modeswitch/pkg/ssh_client"
	"github.com/streamdp/modeswitch/pkg/telnet_client"
)

type Sender interface {
	Send(cmd string) ([]byte, error)
	Close() error
}

type Response struct {
	Text string `json:"text"`
	SMS  *SMSData
}

type cmdFn func(d *Device) (*Response, error)

type Device struct {
	sender Sender
	cfg    *config.UserConfig
}

func New(cfg *config.UserConfig) (*Device, error) {
	var (
		sender Sender
		err    error
	)

	if cfg.IsSsh {
		if sender, err = ssh_client.New(cfg); err != nil {
			return nil, err
		}
	} else {
		if sender, err = telnet_client.New(cfg); err != nil {
			return nil, err
		}
	}

	return &Device{
		sender: sender,
		cfg:    cfg,
	}, nil
}

func (d *Device) Close() error {
	return d.sender.Close()
}

func SendCmd(cfg *config.UserConfig, f cmdFn) (*Response, error) {
	d, err := New(cfg)
	if err != nil {
		return nil, err
	}
	defer func(d *Device) {
		_ = d.Close()
	}(d)

	return f(d)
}
