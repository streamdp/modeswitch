package connections

import (
	"fmt"
	"regexp"
	"time"

	"github.com/streamdp/modeswitch/config"
	"github.com/streamdp/telnet-client"
)

func SendTelnetCommand(mode string, c *config.UserConfig) (err error) {
	cli := telnet.TelnetClient{
		Address:  c.Host,
		Port:     c.Port,
		Login:    c.UserName,
		Password: c.Password,
		Timeout:  5 * time.Second,

		LoginRe:  regexp.MustCompile("Login:"),
		BannerRe: regexp.MustCompile("\\(config\\)>"),
	}

	if err = cli.Dial(); err != nil {
		err = fmt.Errorf("unable to connect:\n%w", err)
		return
	}
	defer cli.Close()

	var init = c.InitLte
	if mode == config.UMTS {
		init = c.InitUmts
	}

	if _, err = cli.Execute(fmt.Sprintf("interface %s lte init %s", c.InterfaceId, init)); err != nil {
		return
	}
	if _, err = cli.Execute(fmt.Sprintf("interface %s usb acq %s", c.InterfaceId, mode)); err != nil {
		return
	}

	return
}
