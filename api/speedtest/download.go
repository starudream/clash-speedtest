package speedtest

import (
	"github.com/starudream/clash-speedtest/api/common"
)

// size see: https://github.com/sivel/speedtest-cli/blob/v2.1.3/speedtest.py#L1186

const size = 1000 // ~ 2 MB

func (c *Client) Download(s *Server, fn common.DownloadBodyFunc) (*common.DownloadResult, error) {
	return common.Download(c.c.R(), s.BaseURL("random%dx%d.jpg", size, size), 0, fn)
}
