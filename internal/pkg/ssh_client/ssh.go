package ssh_client

import (
	"net"

	"github.com/helloyi/go-sshclient"
	"github.com/streamdp/modeswitch/internal/config"
	"golang.org/x/crypto/ssh"
)

type client struct {
	cli *sshclient.Client
}

func New(c *config.UserConfig) (*client, error) {
	cli, err := sshclient.Dial("tcp", net.JoinHostPort(c.Host, c.Port), &ssh.ClientConfig{
		User: c.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
		Timeout:         config.DefaultTimeout,
	})
	if err != nil {
		return nil, err
	}

	return &client{
		cli: cli,
	}, nil
}

func (c *client) Send(cmd string) ([]byte, error) {
	out, err := c.cli.Cmd(cmd).Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *client) Close() error {
	return c.cli.Close()
}
