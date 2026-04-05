package device

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/streamdp/modeswitch/config"
)

var errWrongInterface = errors.New("wrong interface")

func SwitchLTEModemMode(sender Sender, mode string, c *config.UserConfig) error {
	var init = c.InitLte
	if mode == config.UMTS {
		init = c.InitUmts
	}
	output, err := sender.Send(
		fmt.Sprintf("interface %s lte init %s", c.InterfaceId, init),
		fmt.Sprintf("interface %s usb acq %s", c.InterfaceId, mode),
	)
	if err != nil {
		return err
	}

	if bytes.Contains(output, []byte("Interface error")) {
		return errWrongInterface
	}

	return nil
}
