package clash

import (
	"fmt"
	"sort"
)

type GetProvidersResp struct {
	Providers map[string]*Provider `json:"providers"`
}

func (t *GetProvidersResp) FilterProxies(adapters ...string) []*Proxy {
	if len(adapters) == 0 {
		adapters = AvailAdapters
	}
	avail := map[string]bool{}
	for _, adapter := range adapters {
		avail[adapter] = true
	}
	pm := map[string]*Proxy{}
	for _, provider := range t.Providers {
		for _, proxy := range provider.Proxies {
			if avail[proxy.Type] {
				pm[proxy.Name] = proxy
			}
		}
	}
	var proxies []*Proxy
	for _, proxy := range pm {
		proxies = append(proxies, proxy)
	}
	sort.Slice(proxies, func(i, j int) bool { return proxies[i].Name < proxies[j].Name })
	return proxies
}

type Provider struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Proxies []*Proxy `json:"proxies"`
}

func (c *Client) GetProviderProxies() (*GetProvidersResp, error) {
	resp, err := c.R().SetResult(&GetProvidersResp{}).Get(c.Addr + "/providers/proxies")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("response status %s", resp.Status())
	}

	return resp.Result().(*GetProvidersResp), nil
}
