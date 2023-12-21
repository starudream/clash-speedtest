package cloudflare

import (
	"strconv"

	"github.com/starudream/clash-speedtest/api/common"
)

const size int64 = 2000000

func (c *Client) Download(fn common.DownloadBodyFunc) (*common.DownloadResult, error) {
	return common.Download(c.c.R(), "https://speed.cloudflare.com/__down?bytes="+strconv.FormatInt(size, 10), size, fn)
}
