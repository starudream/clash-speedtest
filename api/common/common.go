package common

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/resty/v2"

	"github.com/starudream/clash-speedtest/util"
)

type DownloadResult struct {
	TotalSize int64
	ConnTime  time.Duration
	RespTime  time.Duration
}

func (t *DownloadResult) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TotalSize: %s, ConnTime: %s, RespTime: %s", util.Bytes(t.TotalSize), t.ConnTime, t.RespTime)
}

type DownloadBodyFunc func(body io.ReadCloser, size int64) error

func Download(req *resty.Request, url string, size int64, fn DownloadBodyFunc) (*DownloadResult, error) {
	resp, err := req.EnableTrace().SetDoNotParseResponse(true).Get(url)
	if err != nil {
		return nil, err
	}
	if size <= 0 {
		size, _ = strconv.ParseInt(resp.Header().Get("Content-Length"), 10, 64)
	}
	defer gh.Close(resp.RawBody())
	start := time.Now()
	if fn == nil {
		_, err = io.Copy(io.Discard, resp.RawBody())
	} else {
		err = fn(resp.RawBody(), size)
	}
	if err != nil {
		return nil, err
	}
	res := &DownloadResult{
		TotalSize: size,
		ConnTime:  resp.Request.TraceInfo().ConnTime,
		RespTime:  time.Since(start),
	}
	return res, nil
}
