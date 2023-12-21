package cloudflare

type Config struct {
	Ip      string `header:"Cf-Meta-Ip"`
	Country string `header:"Cf-Meta-Country"`
	Lat     string `header:"Cf-Meta-Latitude"`
	Lon     string `header:"Cf-Meta-Longitude"`
}

func (c *Client) GetConfig() (*Config, error) {
	resp, err := c.c.R().Get("https://speed.cloudflare.com/__down?bytes=1")
	if err != nil {
		return nil, err
	}

	config := &Config{
		Ip:      resp.Header().Get("Cf-Meta-Ip"),
		Country: resp.Header().Get("Cf-Meta-Country"),
		Lat:     resp.Header().Get("Cf-Meta-Latitude"),
		Lon:     resp.Header().Get("Cf-Meta-Longitude"),
	}

	return config, nil
}
