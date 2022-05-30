package download

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	cfMetaCountry = "cf-meta-country"
	cfMetaCity    = "cf-meta-city"
	cfMetaIP      = "cf-meta-ip"
	cfMetaLat     = "cf-meta-latitude"
	cfMetaLng     = "cf-meta-longitude"
)

var cfURL = "https://speed.cloudflare.com/__down?bytes=" + strconv.Itoa(100*1000*1000)

func (c *Client) Cloudflare() (*Result, error) {
	result, err := Do(c.proxy, cfURL, c.timeout)
	if err != nil {
		return nil, err
	}

	if result.status != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", result.status)
	}

	result.IP = result.header.Get(cfMetaIP)
	result.Country = result.header.Get(cfMetaCountry)
	result.City = result.header.Get(cfMetaCity)
	result.Lat = result.header.Get(cfMetaLat)
	result.Lng = result.header.Get(cfMetaLng)

	return result, nil
}
