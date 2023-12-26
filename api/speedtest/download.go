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

const size = 1000 // ~ 2 MB

func (c *Client) Download(s *Server, fn common.DownloadBodyFunc) (*common.DownloadResult, error) {
	return common.Download(c.c.R(), s.BaseURL("random%dx%d.jpg", size, size), 0, fn)
}
