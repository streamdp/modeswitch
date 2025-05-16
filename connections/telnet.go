package connections

import (
	"fmt"
	"regexp"

	"github.com/streamdp/modeswitch/config"
	"github.com/streamdp/telnet-client"
)

func SendTelnetCommand(mode string, c *config.UserConfig) error {
	cli := telnet.TelnetClient{
		Address:  c.Host,
		Port:     c.Port,
		Login:    c.UserName,
		Password: c.Password,

		ConnTimeout: config.DefaultTimeout,
		ReadTimeout: config.DefaultTimeout,

		LoginRe:  regexp.MustCompile("Login:"),
		BannerRe: regexp.MustCompile("\\(config\\)>"),
	}

	if err := cli.Dial(); err != nil {
		return err
	}
	defer cli.Close()

	var init = c.InitLte
	if mode == config.UMTS {
		init = c.InitUmts
	}

	if _, err := cli.Execute(fmt.Sprintf("interface %s lte init %s", c.InterfaceId, init)); err != nil {
		return nil
	}
	if _, err := cli.Execute(fmt.Sprintf("interface %s usb acq %s", c.InterfaceId, mode)); err != nil {
		return nil
	}

	return nil
}
