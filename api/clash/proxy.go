package clash

import (
	"fmt"
	"sort"

	"github.com/starudream/go-lib/core/v2/gh"
)

type GetProxiesResp struct {
	Proxies map[string]*Proxy `json:"proxies"`
}

func (t *GetProxiesResp) FilterProxies(adapters ...string) []*Proxy {
	if len(adapters) == 0 {
		adapters = AvailAdapters
	}
	avail := map[string]bool{}
	for _, adapter := range adapters {
		avail[adapter] = true
	}
	var proxies []*Proxy
	for _, proxy := range t.Proxies {
		if avail[proxy.Type] {
			proxies = append(proxies, proxy)
		}
	}
	sort.Slice(proxies, func(i, j int) bool { return proxies[i].Name < proxies[j].Name })
	return proxies
}

type Proxy struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

const (
	ProxyGlobal = "GLOBAL"
	ProxyDirect = "DIRECT"
)

// AvailAdapters https://github.com/MetaCubeX/mihomo/blob/v1.17.0/constant/adapters.go#L170
var AvailAdapters = []string{
	"Shadowsocks",
	"ShadowsocksR",
	"Snell",
	"Socks5",
	"Http",
	"Vmess",
	"Vless",
	"Trojan",
	"Hysteria",
	"Hysteria2",
	"WireGuard",
	"Tuic",
}

func (c *Client) GetProxies() (*GetProxiesResp, error) {
	resp, err := c.R().SetResult(&GetProxiesResp{}).Get(c.Addr + "/proxies")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("response status %s", resp.Status())
	}

	return resp.Result().(*GetProxiesResp), nil
}

func (c *Client) SetGlobalProxy(name string) error {
	resp, err := c.R().SetBody(gh.M{"name": name}).Put(c.Addr + "/proxies/" + ProxyGlobal)
	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("response status %s", resp.Status())
	}

	return nil
}
