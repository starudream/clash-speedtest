package clash

import (
	"strings"
	"time"

	"github.com/starudream/go-lib/resty/v2"
)

type Client struct {
	Addr   string
	Secret string

	c *resty.Client
}

func NewClient(addr, secret string) *Client {
	return &Client{
		Addr:   strings.TrimSuffix(addr, "/"),
		Secret: secret,
		c:      resty.New().SetTimeout(time.Minute),
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.c.SetTimeout(timeout)
	return c
}

func (c *Client) R() *resty.Request {
	r := c.c.R()
	if c.Secret != "" {
		r.SetHeader("Authorization", "Bearer "+c.Secret)
	}
	return r
}
