package speedtest

import (
	"time"

	"github.com/starudream/go-lib/resty/v2"
)

const Name = "speedtest"

type Client struct {
	c *resty.Client
}

func NewClient() *Client {
	return &Client{resty.New().SetTimeout(time.Minute)}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.c.SetTimeout(timeout)
	return c
}

func (c *Client) WithProxy(proxy string) *Client {
	c.c.SetProxy(proxy)
	return c
}
