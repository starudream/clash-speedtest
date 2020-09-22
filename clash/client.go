package clash

import (
	"fmt"

	"github.com/starudream/clash-speedtest/util"
)

type Client struct {
	url    string
	secret string
}

func New() *Client {
	return &Client{}
}

func (c *Client) SetURL(u string) *Client {
	c.url = u
	return c
}

func (c *Client) SetSecret(s string) *Client {
	c.secret = s
	return c
}

func (c *Client) headers() map[string]string {
	headers := map[string]string{}
	if c.secret != "" {
		headers["Authorization"] = "Bearer " + c.secret
	}
	if len(headers) == 0 {
		return nil
	}
	return headers
}

type Version struct {
	Version string `json:"version"`
	Premium bool   `json:"premium"`
}

func (c *Client) GetVersion() (version *Version, err error) {
	_, err = util.HTTPGet(c.url+"/version", c.headers(), &version)
	return
}

type Config struct {
	Mode string `json:"mode"`
}

type Mode string

const (
	ModeNone   Mode = ""
	ModeGlobal Mode = "global"
	ModeRule   Mode = "rule"
	ModeDirect Mode = "direct"
)

var ModeMap = map[string]Mode{
	string(ModeNone):   ModeNone,
	string(ModeGlobal): ModeGlobal,
	string(ModeRule):   ModeRule,
	string(ModeDirect): ModeDirect,
}

func (mode Mode) String() string {
	return string(mode)
}

func (c *Client) GetConfigMode() (Mode, error) {
	var config *Config
	_, err := util.HTTPGet(c.url+"/configs", c.headers(), &config)
	if err != nil {
		return "", err
	}
	mode, ok := ModeMap[config.Mode]
	if ok {
		return mode, nil
	}
	return ModeNone, fmt.Errorf("clash: %s not found", config.Mode)
}

func (c *Client) PatchConfigMode(mode Mode) (err error) {
	_, err = util.HTTPPatch(c.url+"/configs", c.headers(), map[string]string{"mode": mode.String()}, nil)
	return
}

type Proxies struct {
	Proxies map[string]Proxy `json:"proxies"`
}

type Proxy struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (c *Client) GetProxies() (proxies *Proxies, err error) {
	_, err = util.HTTPGet(c.url+"/proxies", c.headers(), &proxies)
	return
}

func (c *Client) PutProxiesGlobal(name string) (err error) {
	_, err = util.HTTPPut(c.url+"/proxies/GLOBAL", c.headers(), map[string]string{"name": name}, nil)
	return
}
