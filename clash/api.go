package clash

import (
	"fmt"
	"net/http"
	"strings"
)

type VersionResp struct {
	Premium bool   `json:"premium"`
	Version string `json:"version"`
}

func (c *Client) GetVersion() (*VersionResp, error) {
	return do[*VersionResp](c, http.MethodGet, "/version", nil)
}

type GetModeResp struct {
	Mode string `json:"mode"`
}

func (r *GetModeResp) GetMode() Mode {
	if r == nil {
		return ""
	}
	return ModeMap[strings.ToLower(r.Mode)]
}

func (c *Client) GetMode() (Mode, error) {
	resp, err := do[*GetModeResp](c, http.MethodGet, "/configs", nil)
	return resp.GetMode(), err
}

func (c *Client) SetMode(mode Mode) error {
	_, err := do[any](c, http.MethodPatch, "/configs", fmt.Sprintf(`{"mode":"%s"}`, mode))
	return err
}

type GetProxiesResp struct {
	Proxies map[string]*Proxy `json:"proxies"`
}

func (r *GetProxiesResp) GetProxies() map[string]*Proxy {
	if r == nil {
		return nil
	}
	return r.Proxies
}

type Proxy struct {
	Name string   `json:"name"`
	Type string   `json:"type"`
	All  []string `json:"all"`
	Now  string   `json:"now"`
}

func (c *Client) GetProxies() (map[string]*Proxy, error) {
	resp, err := do[*GetProxiesResp](c, http.MethodGet, "/proxies", nil)
	return resp.GetProxies(), err
}

func (c *Client) SetProxy(name string) error {
	/// TODO: 改成配置文件下发的手动规则组
	_, err := do[any](c, http.MethodPut, "/proxies/🚀 手动切换", fmt.Sprintf(`{"name":"%s"}`, name))
	return err
}
