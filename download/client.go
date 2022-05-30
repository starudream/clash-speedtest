package download

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

type Client struct {
	proxy   string
	timeout time.Duration
}

func New(proxy string, timeout time.Duration) *Client {
	return &Client{proxy: proxy, timeout: timeout}
}

type Result struct {
	BeginAt         time.Time
	ConnectDuration time.Duration
	TotalDuration   time.Duration

	Length int64

	BS float64

	IP      string
	Country string
	City    string
	Lat     string
	Lng     string

	Name string
	Type string

	status int
	header http.Header
}

func Do(proxy, url string, timeout time.Duration) (*Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	result := &Result{BeginAt: time.Now()}

	resp, err := newClient(proxy).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result.ConnectDuration = time.Since(result.BeginAt)

	result.status = resp.StatusCode
	result.header = resp.Header

	flag := int64(0)

	go func() {
		for {
			bs := make([]byte, 100*1000)
			n, re := resp.Body.Read(bs)
			if re != nil {
				cancel()
				return
			}
			atomic.AddInt64(&flag, int64(n))
		}
	}()

	<-ctx.Done()

	result.TotalDuration = time.Since(result.BeginAt)

	result.Length = atomic.LoadInt64(&flag)

	result.BS = float64(result.Length) / result.TotalDuration.Seconds()

	return result, nil
}

func newClient(proxy string) *http.Client {
	var proxyFunc func(*http.Request) (*url.URL, error)
	if proxy != "" {
		proxyFunc = func(req *http.Request) (*url.URL, error) {
			return url.Parse(proxy)
		}
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy:               proxyFunc,
			DialContext:         (&net.Dialer{Timeout: time.Hour}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
			IdleConnTimeout:     time.Minute,
			ForceAttemptHTTP2:   true,
		},
		Timeout: time.Hour,
	}
}
