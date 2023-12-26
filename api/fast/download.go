package fast

import (
	"net/url"
	"strconv"

	"github.com/starudream/clash-speedtest/api/common"
)

const (
	miniSize    = 1 * 1e6
	defaultSize = 4 * 1e6
)

func (c *Client) Download(t *Target, size int, fn common.DownloadBodyFunc) (*common.DownloadResult, error) {
	u, err := url.Parse(t.Url)
	if err != nil {
		return nil, err
	}
	size *= 1e6
	if size <= miniSize {
		size = defaultSize
	}
	u.Path += "/range/0-" + strconv.FormatInt(int64(size), 10)
	return common.Download(c.c.R(), u.String(), int64(size), fn)
}
