package device

import (
	"fmt"

	"github.com/streamdp/modeswitch/config"
)

func (d *Device) SwitchLTEModemMode(mode string) (*Response, error) {
	var init = d.cfg.InitLte
	if mode == config.UMTS {
		init = d.cfg.InitUmts
	}
	out, err := d.sender.Send(
		fmt.Sprintf("interface %s lte init %s", d.cfg.InterfaceId, init),
		fmt.Sprintf("interface %s usb acq %s", d.cfg.InterfaceId, mode),
	)
	if err != nil {
		return nil, err
	}

	if err = errorLookup(out); err != nil {
		return nil, err
	}

	return &Response{Text: string(out)}, nil
}
