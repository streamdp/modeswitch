package device

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/streamdp/modeswitch/config"
)

var errWrongInterface = errors.New("wrong interface")

func (d *Device) SwitchLTEModemMode(mode string) (*Response, error) {
	var init = d.cfg.InitLte
	if mode == config.UMTS {
		init = d.cfg.InitUmts
	}
	output, err := d.sender.Send(
		fmt.Sprintf("interface %s lte init %s", d.cfg.InterfaceId, init),
		fmt.Sprintf("interface %s usb acq %s", d.cfg.InterfaceId, mode),
	)
	if err != nil {
		return nil, err
	}

	if bytes.Contains(output, []byte("Interface error")) {
		return nil, errWrongInterface
	}

	return &Response{Text: string(output)}, nil
}
