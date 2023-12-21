package clash

import (
	"fmt"

	"github.com/starudream/go-lib/core/v2/gh"
)

type Config struct {
	Port      int `json:"port"`
	SocksPort int `json:"socks-port"`
	MixedPort int `json:"mixed-port"`

	Authentication []string `json:"authentication"`

	Mode Mode `json:"mode"`
}

//goland:noinspection HttpUrlsUsage
func (c *Config) Proxy(host string) (string, error) {
	if len(c.Authentication) > 0 {
		return "", fmt.Errorf("authentication is enabled, please set the clash-proxy manually")
	}
	if c.Port > 0 {
		return fmt.Sprintf("http://%s:%d", host, c.Port), nil
	} else if c.SocksPort > 0 {
		return fmt.Sprintf("socks5://%s:%d", host, c.SocksPort), nil
	} else if c.MixedPort > 0 {
		return fmt.Sprintf("http://%s:%d", host, c.MixedPort), nil
	}
	return "", fmt.Errorf("no detected proxy port")
}

type Mode string

const (
	ModeGlobal Mode = "global"
	ModeRule   Mode = "rule"
	ModeDirect Mode = "direct"
)

func (c *Client) GetConfig() (*Config, error) {
	resp, err := c.R().SetResult(&Config{}).Get(c.Addr + "/configs")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("response status %s", resp.Status())
	}

	return resp.Result().(*Config), nil
}

func (c *Client) SetMode(mode Mode) error {
	resp, err := c.R().SetBody(gh.M{"mode": mode}).Patch(c.Addr + "/configs")
	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("response status %s", resp.Status())
	}

	return nil
}
