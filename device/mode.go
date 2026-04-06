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

	var output []byte
	for _, cmd := range []string{
		fmt.Sprintf("interface %s lte init %s", d.cfg.InterfaceId, init),
		fmt.Sprintf("interface %s usb acq %s", d.cfg.InterfaceId, mode),
		"system configuration save",
		"exit",
	} {
		out, err := d.sender.Send(cmd)
		if err != nil {
			return nil, err
		}

		if err = errorLookup(out); err != nil {
			return nil, err
		}

		if out != nil {
			output = append(output, out...)
		}
	}

	return &Response{Text: string(output)}, nil
}
