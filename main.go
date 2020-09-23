package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-sdk/logx"
	"github.com/go-sdk/utilx/json"

	"github.com/starudream/clash-speedtest/clash"
	"github.com/starudream/clash-speedtest/fast"
	"github.com/starudream/clash-speedtest/util"
)

type Config struct {
	URL     string `json:"url"`
	Secret  string `json:"secret"`
	Proxy   string `json:"proxy"`
	Include string `json:"include,omitempty"`
	Exclude string `json:"exclude,omitempty"`
	Process bool   `json:"process"`
	Help    bool   `json:"-"`
}

type Dashboard struct {
	TotalBytes int64         `json:"total_bytes"`
	TotalTime  time.Duration `json:"total_time"`
	Nodes      []*Node       `json:"nodes"`
}

type Node struct {
	Name  string `json:"name"`
	Speed string `json:"speed"` // kb/s
}

const (
	MaxRetry = 3
)

var (
	config = &Config{}

	client *clash.Client
)

func init() {
	flag.StringVar(&config.URL, "url", "http://127.0.0.1:9090", "external controller url")
	flag.StringVar(&config.Secret, "secret", "", "external controller secret")
	flag.StringVar(&config.Proxy, "proxy", "http://127.0.0.1:7890", "http proxy url")
	flag.StringVar(&config.Include, "include", "", "filter nodes that include")
	flag.StringVar(&config.Exclude, "exclude", "", "filter nodes that exclude")
	flag.BoolVar(&config.Process, "process", false, "show speedtest process")
	flag.BoolVar(&config.Help, "help", false, "instructions for use")
	flag.Parse()

	if config.Help {
		flag.Usage()
		os.Exit(0)
	}

	logx.Infof("[config] %s", json.MustMarshal(config))

	if config.URL == "" {
		logx.Fatal("[config] external controller url is empty")
	}
	if config.Proxy == "" {
		logx.Fatal("[config] http proxy url is empty")
	}

	client = clash.New().SetURL(config.URL).SetSecret(config.Secret)

	version, err := client.GetVersion()
	if err != nil {
		logx.WithField("err", err).Fatal("[clash] get version fail")
	}

	logx.Infof("[clash] %s", json.MustMarshal(version))
}

func main() {
	hp, hsp := util.ProxyGet()

	util.ProxySet("", "")

	mode, err := client.GetConfigMode()
	if err != nil {
		logx.WithField("err", err).Fatal("[clash] get proxy mode fail")
	}

	err = client.PatchConfigMode(clash.ModeGlobal)
	if err != nil {
		logx.WithField("err", err).Fatal("[clash] switch mode to GLOBAL fail")
	}

	logx.Info("[clash] switch mode to GLOBAL success")

	defer func() {
		util.ProxySet(hp, hsp)
		if mode != clash.ModeGlobal {
			err := client.PatchConfigMode(mode)
			if err != nil {
				logx.WithField("err", err).Fatalf("[clash] recovery mode to %s fail, please switch manually", strings.ToUpper(mode.String()))
			}
			logx.Infof("[clash] recovery mode to %s success", strings.ToUpper(mode.String()))
		}
	}()

	proxies, err := client.GetProxies()
	if err != nil {
		logx.WithField("err", err).Fatal("[clash] get proxies fail")
	}

	var names []string
	for _, proxy := range proxies.Proxies {
		switch proxy.Type {
		case "Shadowsocks", "Vmess":
		default:
			continue
		}
		if config.Include != "" && !strings.Contains(proxy.Name, config.Include) {
			continue
		}
		if config.Exclude != "" && strings.Contains(proxy.Name, config.Exclude) {
			continue
		}
		names = append(names, proxy.Name)
	}
	sort.Strings(names)

	if len(names) == 0 {
		logx.Fatal("[config] no nodes left, please change include and exclude arguments")
	}

	nameMaxLen := 30

	logx.Infof("[speedtest] total nodes: %d", len(names))
	for i := 0; i < len(names); i++ {
		name := names[i]
		logx.Infof("-> %s", name)
		nameLen := len(name)
		if nameMaxLen < nameLen {
			nameMaxLen = nameLen
		}
	}

	retry := func(i int) {
		time.Sleep(time.Duration(i) * time.Second)
		logx.Warnf("[speedtest] attempts %d time(s)", i)
	}

	dashboard := &Dashboard{Nodes: make([]*Node, len(names))}

	for i := 0; i < len(names); i++ {
		proxy := proxies.Proxies[names[i]]

		err := client.PutProxiesGlobal(proxy.Name)
		if err != nil {
			logx.WithField("err", err).Fatalf("[clash] switch node fail")
		}

		time.Sleep(time.Second)

		util.ProxySet(config.Proxy, config.Proxy)

		for j := 1; j <= MaxRetry; j++ {
			data, err := fast.GetData()
			if err != nil {
				logx.WithField("err", err).Errorf("[fast.com] api fail")
				retry(j)
				continue
			}

			logx.Infof("[%s] (%s) country: %s, city: %s", proxy.Name, data.Client.IP, data.Client.Location.Country, data.Client.Location.City)

			if len(data.Targets) == 0 {
				logx.Errorf("[%s] current area not exist speedtest node", proxy.Name)
				break
			}

			target := data.Targets[0]

			logx.Infof("[%s] speedtest node country: %s, city: %s", proxy.Name, target.Location.Country, target.Location.City)

			result, err := SpeedTest(target.URL, 0, config.Process)
			if err != nil {
				logx.WithField("err", err).Errorf("[%s] speedtest fail", proxy.Name)
				continue
			}

			kb := float64(result.TotalBytes) / 1024
			ti := float64(result.TotalTime) / float64(time.Second)
			logx.Infof("[%s] speedtest download: %d kb, took: %.03f s, speed: %.2f kb/s", proxy.Name, int64(kb), ti, kb/ti)

			dashboard.TotalBytes += result.TotalBytes
			dashboard.TotalTime += result.TotalTime
			dashboard.Nodes[i] = &Node{Name: proxy.Name, Speed: fmt.Sprintf("%.2f", kb/ti)}

			logx.Infof("[%s] speedtest done, %d/%d", proxy.Name, i+1, len(names))
			break
		}

		util.ProxySet("", "")
	}

	logx.Infof("total bytes: %.02f mb, total time: %d s", float64(dashboard.TotalBytes)/1024/1024, int64(dashboard.TotalTime/time.Second))

	format := fmt.Sprintf("-> %%%ds   %%15s", nameMaxLen)

	logx.Infof(format, "name", "speed (kb/s)")

	for i := 0; i < len(names); i++ {
		logx.Infof(format, names[i], dashboard.Nodes[i].Speed)
	}
}
