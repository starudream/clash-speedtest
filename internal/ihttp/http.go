package ihttp

import (
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	hdrUserAgentKey   = http.CanonicalHeaderKey("User-Agent")
	hdrUserAgentValue = "Golang/" + strings.TrimLeft(runtime.Version(), "go")
)

func New() *resty.Client {
	nlog := &logger{log.Logger}

	c := resty.New()
	c.SetTimeout(5 * time.Minute)
	c.SetLogger(nlog)
	c.SetDisableWarn(true)
	c.SetDebug(viper.GetBool("debug"))
	c.OnBeforeRequest(func(_ *resty.Client, req *resty.Request) error {
		if strings.TrimSpace(req.Header.Get(hdrUserAgentKey)) == "" {
			req.Header.Set(hdrUserAgentKey, hdrUserAgentValue)
		}
		return nil
	})
	return c
}

var (
	client *resty.Client

	clientOnce sync.Once
)

func R() *resty.Request {
	clientOnce.Do(func() {
		client = New()
	})
	return client.R()
}
