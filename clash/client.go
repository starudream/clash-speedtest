package clash

import (
	"fmt"

	"github.com/starudream/go-lib/httpx"
)

const (
	hdrAuthKey    = "Authorization"
	hdrAuthBearer = "Bearer "
)

type Client struct {
	url    string
	secret string

	headers map[string]string
}

func New(url, secret string) *Client {
	c := &Client{url: url, secret: secret, headers: map[string]string{}}
	if secret != "" {
		c.headers = map[string]string{hdrAuthKey: hdrAuthBearer + secret}
	}
	return c
}

type CommonResp struct {
	Message string `json:"message"`
}

func (r *CommonResp) GetMessage() string {
	if r == nil {
		return ""
	}
	return r.Message
}

func do[T any](c *Client, method, path string, body any) (T, error) {
	var result T
	resp, err := httpx.R().SetHeaders(c.headers).SetBody(body).SetResult(&result).SetError(&CommonResp{}).Execute(method, c.url+path)
	if err != nil {
		return result, err
	}
	if resp.IsError() {
		return result, fmt.Errorf("api response error: %s %s", resp.Status(), resp.Error().(*CommonResp).GetMessage())
	}
	return result, err
}
