package job

import (
	"fmt"
	"io"
	"sync"

	"github.com/starudream/go-lib/core/v2/slog"

	"github.com/starudream/clash-speedtest/api/cloudflare"
	"github.com/starudream/clash-speedtest/api/common"
	"github.com/starudream/clash-speedtest/api/fast"
	"github.com/starudream/clash-speedtest/api/speedtest"
	"github.com/starudream/clash-speedtest/util"
)

func (t *Task) Down(result *Result) (err error) {
	var down downFunc
	switch t.Download {
	case cloudflare.Name:
		down, err = t.downCloudflare(result)
	case speedtest.Name:
		down, err = t.downSpeedtest(result)
	case fast.Name:
		down, err = t.downFast(result)
	default:
		return fmt.Errorf("unknown download type: %s", t.Download)
	}
	if err != nil {
		return err
	}

	pool, bars := util.NewBarsPool(t.Threads, result.Proxy.Name)

	wg := sync.WaitGroup{}
	wg.Add(t.Threads)

	for i := 0; i < t.Threads; i++ {
		go func(i int) { down(i, bars[i], &wg) }(i)
	}

	wg.Wait()

	_ = pool.Stop()

	return nil
}

type downFunc func(i int, bar *util.ProgressBar, wg *sync.WaitGroup)

func (t *Task) downCloudflare(result *Result) (downFunc, error) {
	cli := cloudflare.NewClient().WithProxy(t.ClashProxy)

	cfg, err := cli.GetConfig()
	if err != nil {
		return nil, err
	}
	result.Ip = cfg.Ip
	result.Country = cfg.Country
	result.Lat = cfg.Lat
	result.Lon = cfg.Lon

	fn := func(i int, bar *util.ProgressBar, wg *sync.WaitGroup) {
		defer wg.Done()
		res, err2 := cli.Download(t.Size, t.progress(bar))
		if err2 != nil {
			slog.Error("cloudflare error: %v", err2)
			return
		}
		result.SetDownload(i, res)
	}

	return fn, nil
}

func (t *Task) downSpeedtest(result *Result) (downFunc, error) {
	cli := speedtest.NewClient().WithProxy(t.ClashProxy)

	cfg, err := cli.GetConfig()
	if err != nil {
		return nil, err
	}
	result.Ip = cfg.Client.Ip
	result.Country = cfg.Client.Country
	result.Lat = cfg.Client.Lat
	result.Lon = cfg.Client.Lon

	servers, err := cli.GetServers()
	if err != nil {
		return nil, err
	}

	fn := func(i int, bar *util.ProgressBar, wg *sync.WaitGroup) {
		defer wg.Done()
		server := servers[i%len(servers)]
		res, err2 := cli.Download(server, t.Size, t.progress(bar))
		if err2 != nil {
			slog.Error("speedtest error: %v", err2)
			return
		}
		result.SetDownload(i, res)
	}

	return fn, nil
}

func (t *Task) downFast(result *Result) (downFunc, error) {
	cli := fast.NewClient().WithProxy(t.ClashProxy)

	cfg, err := cli.GetConfig()
	if err != nil {
		return nil, err
	}
	result.Ip = cfg.Client.Ip
	result.Country = cfg.Client.Location.Country
	result.City = cfg.Client.Location.City

	fn := func(i int, bar *util.ProgressBar, wg *sync.WaitGroup) {
		defer wg.Done()
		target := cfg.Targets[i%len(cfg.Targets)]
		res, err2 := cli.Download(target, t.Size, t.progress(bar))
		if err2 != nil {
			slog.Error("fast error: %v", err2)
			return
		}
		result.SetDownload(i, res)
	}

	return fn, nil
}

func (t *Task) progress(bar *util.ProgressBar) common.DownloadBodyFunc {
	return func(body io.ReadCloser, size int64) error {
		defer bar.Finish()
		bar.SetTotal(size)
		_, err := io.ReadAll(bar.NewProxyReader(body))
		return err
	}
}
