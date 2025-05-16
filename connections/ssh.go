package connections

import (
	"fmt"
	"net"

	"github.com/helloyi/go-sshclient"
	"github.com/streamdp/modeswitch/config"
	"golang.org/x/crypto/ssh"
)

func SendSshCommand(mode string, c *config.UserConfig) error {
	cli, err := sshclient.Dial("tcp", c.Host+":"+c.Port, &ssh.ClientConfig{
		User: c.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
		Timeout:         config.DefaultTimeout,
	})
	if err != nil {
		return err
	}
	defer func(cli *sshclient.Client) {
		_ = cli.Close()
	}(cli)

	var init = c.InitLte
	if mode == config.UMTS {
		init = c.InitUmts
	}

	if err = cli.Cmd(fmt.Sprintf("interface %s lte init %s", c.InterfaceId, init)).
		Cmd(fmt.Sprintf("interface %s usb acq %s", c.InterfaceId, mode)).Run(); err != nil {
		return err
	}

	return nil
}
