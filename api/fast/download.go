package fast

import (
	"net/url"
	"strconv"

	"github.com/starudream/clash-speedtest/api/common"
)

const size int64 = 2000000

func (c *Client) Download(t *Target, fn common.DownloadBodyFunc) (*common.DownloadResult, error) {
	u, err := url.Parse(t.Url)
	if err != nil {
		return nil, err
	}
	u.Path += "/range/0-" + strconv.FormatInt(size, 10)
	return common.Download(c.c.R(), u.String(), size, fn)
}
