package job

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/starudream/go-lib/core/v2/config"
	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/core/v2/utils/fmtutil"
	"github.com/starudream/go-lib/core/v2/utils/maputil"
	"github.com/starudream/go-lib/tablew/v2"

	"github.com/starudream/clash-speedtest/api/clash"
	"github.com/starudream/clash-speedtest/util"
)

type Task struct {
	ClashAddr   string `yaml:"clash.addr"`
	ClashSecret string `yaml:"clash.secret"`
	ClashProxy  string `yaml:"clash.proxy"`

	Threads  int      `yaml:"threads"`
	Download string   `yaml:"download"`
	Includes []string `yaml:"includes"`
	Excludes []string `yaml:"excludes"`
	Confirm  bool     `yaml:"confirm"`
	Output   string   `yaml:"output"`

	clash   *clash.Client
	version *clash.Version
	config  *clash.Config
	proxies []*clash.Proxy

	results maputil.SyncMap[string, *Result]
}

func Run() error {
	t := &Task{}
	err := config.Unmarshal("", t)
	if err != nil {
		return err
	}

	err = t.Clash()
	if err != nil {
		return err
	}

	if len(t.proxies) == 0 {
		return fmt.Errorf("no proxies found")
	}

	slog.Info("clash version: %s", t.version.Version,
		slog.Bool("premium", t.version.Premium),
		slog.Bool("meta", t.version.Meta),
		slog.String("mode", string(t.config.Mode)),
		slog.String("proxy", t.ClashProxy),
		slog.String("includes", strings.Join(t.Includes, ",")),
		slog.String("excludes", strings.Join(t.Excludes, ",")),
		slog.Int("total", len(t.proxies)),
		slog.Int("threads", t.Threads),
	)

	fmt.Print(tablew.Structs(t.proxies))

	if !t.Confirm {
		input := fmtutil.Scan("confirm to start? [y/n]: ")
		if !strings.EqualFold(input, "y") {
			fmt.Println("cancel by user")
			return nil
		}
	}

	err = t.clash.SetMode(clash.ModeGlobal)
	if err != nil {
		return err
	}
	defer func(start time.Time) {
		err = t.clash.SetMode(t.config.Mode)
		if err != nil {
			slog.Error("set mode error: %v", err)
		}
		slog.Info("took %s", time.Since(start).Truncate(time.Millisecond))
	}(time.Now())

	defer t.Render()

	for i := 0; i < len(t.proxies); i++ {
		proxy := t.proxies[i]

		index := fmt.Sprintf("%d/%d", i+1, len(t.proxies))

		slog.Info("start proxy test: %s", proxy.Name, slog.String("type", proxy.Type), slog.String("index", index))

		result, err2 := t.Test(proxy)
		if err2 != nil {
			slog.Error("test error: %v", err2, slog.String("index", index))
			continue
		}

		slog.Info("proxy test done: %s", proxy.Name, slog.String("index", index))

		t.results.Store(proxy.Name, result)
	}

	return nil
}

func (t *Task) Clash() error {
	u, err := url.Parse(t.ClashAddr)
	if err != nil {
		return err
	}
	h, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return err
	}

	t.clash = clash.NewClient(t.ClashAddr, t.ClashSecret)

	t.version, err = t.clash.GetVersion()
	if err != nil {
		return err
	}

	t.config, err = t.clash.GetConfig()
	if err != nil {
		return err
	}

	if t.ClashProxy == "" {
		t.ClashProxy, err = t.config.Proxy(h)
		if err != nil {
			return err
		}
	}

	providers, err := t.clash.GetProviderProxies()
	if err != nil {
		return err
	}

	for _, proxy := range providers.FilterProxies() {
		if len(t.Includes) > 0 && !util.Contains(t.Includes, proxy.Name) {
			continue
		}
		if len(t.Excludes) > 0 && util.Contains(t.Excludes, proxy.Name) {
			continue
		}
		t.proxies = append(t.proxies, proxy)
	}

	return nil
}

func (t *Task) Render() {
	fi, err := os.Stat(t.Output)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			slog.Error("stat error: %v", err)
			return
		}
		err = os.MkdirAll(t.Output, 0755)
		if err != nil {
			slog.Error("mkdir error: %v", err)
			return
		}
	} else if !fi.IsDir() {
		slog.Error("output is not a directory")
		return
	}

	filename := filepath.Join(t.Output, fmt.Sprintf("%s.txt", time.Now().Format("20060102-150405")))

	table := tablew.Render(func(w *tablew.Table) {
		w.SetAlignment(tablew.ALIGN_CENTER)
		w.SetHeader([]string{"name", "type", "ip", "country", "conn", "down"})
		for i := 0; i < len(t.proxies); i++ {
			proxy := t.proxies[i]
			res, exists := t.results.Load(proxy.Name)
			if !exists {
				continue
			}
			conn := (res.total.ConnTime / time.Duration(res.threads)).Truncate(time.Millisecond)
			down := int64(float64(res.total.TotalSize) / res.total.RespTime.Seconds())
			w.Rich(
				[]string{proxy.Name, proxy.Type, res.Ip, res.Country, conn.String(), util.BytesSec(down)},
				[]tablew.Colors{{tablew.Bold}, {}, {}, {}, {connColor(conn)}, {tablew.Bold, downColor(down)}},
			)
		}
	})
	fmt.Print(table)

	err = os.WriteFile(filename, []byte(table), 0644)
	if err != nil {
		slog.Error("write file error: %v", err)
		return
	}
}

func connColor(d time.Duration) int {
	if d >= time.Second {
		return tablew.FgRedColor
	}
	return 0
}

func downColor(i int64) int {
	if i >= 5*1e6 {
		return tablew.FgGreenColor
	} else if i >= 2*1e6 {
		return tablew.FgYellowColor
	}
	return tablew.FgRedColor
}
