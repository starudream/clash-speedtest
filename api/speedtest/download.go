package speedtest

import (
	"github.com/starudream/clash-speedtest/api/common"
)

// size see: https://github.com/sivel/speedtest-cli/blob/v2.1.3/speedtest.py#L1186

// 1000 ~ 1.99 MB
// 1500 ~ 4.47 MB
// 2000 ~ 7.91 MB
// 2500 ~ 12.41 MB
// 3000 ~ 17.82 MB
// 3500 ~ 24.26 MB
// 4000 ~ 31.63 MB

const defaultSize = 1000 // ~ 2 MB

var availSizes = map[int]bool{
	1000: true,
	1500: true,
	2000: true,
	2500: true,
	3000: true,
	3500: true,
	4000: true,
}

func (c *Client) Download(s *Server, size int, fn common.DownloadBodyFunc) (*common.DownloadResult, error) {
	if !availSizes[size] {
		size = defaultSize
	}
	return common.Download(c.c.R(), s.BaseURL("random%dx%d.jpg", defaultSize, defaultSize), 0, fn)
}
