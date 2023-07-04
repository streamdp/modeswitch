package connections

import (
	"fmt"

	"github.com/helloyi/go-sshclient"
	"github.com/streamdp/modeswitch/config"
)

func SendSshCommand(mode string, c *config.UserConfig) (err error) {
	var cli *sshclient.Client
	if cli, err = sshclient.DialWithPasswd(c.Host+":"+c.Port, c.UserName, c.Password); err != nil {
		return
	}
	defer cli.Close()

	var init = c.InitLte
	if mode == config.UMTS {
		init = c.InitUmts
	}

	if err = cli.Cmd(fmt.Sprintf("interface %s lte init %s", c.InterfaceId, init)).
		Cmd(fmt.Sprintf("interface %s usb acq %s", c.InterfaceId, mode)).Run(); err != nil {
		return
	}

	return
}
