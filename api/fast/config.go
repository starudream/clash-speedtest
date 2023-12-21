package fast

import (
	"fmt"
	"regexp"

	"github.com/starudream/go-lib/core/v2/gh"
)

type Config struct {
	Token  string
	Client struct {
		Ip       string `json:"ip"`
		Asn      string `json:"asn"`
		Location struct {
			Country string `json:"country"`
			City    string `json:"city"`
		} `json:"location"`
	} `json:"client"`
	Targets []*Target `json:"targets"`
}

type Target struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Location struct {
		Country string `json:"country"`
		City    string `json:"city"`
	} `json:"location"`
}

func (c *Client) GetConfig() (*Config, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	resp, err := c.c.R().SetResult(&Config{}).SetQueryParams(gh.MS{"https": "true", "token": token, "urlCount": "5"}).Get("https://api.fast.com/netflix/speedtest/v2")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("response status %s", resp.Status())
	}

	config := resp.Result().(*Config)
	config.Token = token

	return config, nil
}

var (
	scriptRegex = regexp.MustCompile(`(?Us)<script src="(.+)">`)
	tokenRegex  = regexp.MustCompile(`(?U)token:"(.+)"`)
)

func (c *Client) GetToken() (string, error) {
	resp1, err := c.c.R().Get("https://fast.com")
	if err != nil {
		return "", err
	}

	if !resp1.IsSuccess() {
		return "", fmt.Errorf("response status %s", resp1.Status())
	}

	sub1 := scriptRegex.FindStringSubmatch(string(resp1.Body()))
	if len(sub1) != 2 {
		return "", fmt.Errorf("script not found")
	}

	resp2, err := c.c.R().Get("https://fast.com" + sub1[1])
	if err != nil {
		return "", err
	}

	if !resp2.IsSuccess() {
		return "", fmt.Errorf("response status %s", resp2.Status())
	}

	sub2 := tokenRegex.FindStringSubmatch(string(resp2.Body()))
	if len(sub2) != 2 {
		return "", fmt.Errorf("token not found")
	}

	return sub2[1], nil
}
