package telnet_client

import (
	"bytes"
	"regexp"

	"github.com/streamdp/modeswitch/config"
	"github.com/streamdp/telnet-client"
)

type client struct {
	cli *telnet.TelnetClient
}

func New(c *config.UserConfig) (*client, error) {
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
		return nil, err
	}

	return &client{
		cli: &cli,
	}, nil
}

func (c *client) Send(command ...string) ([]byte, error) {
	var output []byte

	for _, cmd := range command {
		out, err := c.cli.Execute(cmd)
		if err != nil {
			return nil, err
		}
		if out != nil {
			output = append(output, out...)
		}
	}

	return bytes.ReplaceAll(output, []byte{27, 91, 75}, nil), nil
}

func (c *client) Close() error {
	c.cli.Close()

	return nil
}
