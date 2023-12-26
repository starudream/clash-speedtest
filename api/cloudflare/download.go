package cloudflare

import (
	"strconv"

	"github.com/starudream/clash-speedtest/api/common"
)

const (
	miniSize    = 1 * 1e6
	defaultSize = 4 * 1e6
)

func (c *Client) Download(size int, fn common.DownloadBodyFunc) (*common.DownloadResult, error) {
	size *= 1e6
	if size <= miniSize {
		size = defaultSize
	}
	url := "https://speed.cloudflare.com/__down?bytes=" + strconv.FormatInt(int64(size), 10)
	return common.Download(c.c.R(), url, int64(size), fn)
}
