package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-sdk/logx"
	"github.com/go-sdk/utilx/json"
	"github.com/olekukonko/tablewriter"

	"github.com/starudream/clash-speedtest/clash"
	"github.com/starudream/clash-speedtest/fast"
	"github.com/starudream/clash-speedtest/util"
)

type Config struct {
	URL     string       `json:"url"`
	Secret  string       `json:"secret"`
	Proxy   string       `json:"proxy"`
	Include stringsValue `json:"include,omitempty"`
	Exclude stringsValue `json:"exclude,omitempty"`
	Retry   int64        `json:"retry"`
	Timeout int64        `json:"timeout"` // seconds
	Process bool         `json:"process"`
	Help    bool         `json:"-"`
}

type Dashboard struct {
	TotalBytes int64         `json:"total_bytes"`
	TotalTime  time.Duration `json:"total_time"`
	Nodes      []*Node       `json:"nodes"`
}

type Node struct {
	Name  string  `json:"name"`
	Speed float64 `json:"speed"` // kb/s
}

var (
	config = &Config{}

	client *clash.Client
)

func init() {
	flag.StringVar(&config.URL, "url", "http://127.0.0.1:9090", "external controller url")
	flag.StringVar(&config.Secret, "secret", "", "external controller secret")
	flag.StringVar(&config.Proxy, "proxy", "http://127.0.0.1:7890", "http proxy url")

	flag.Var(&config.Include, "include", "filter nodes that include")
	flag.Var(&config.Exclude, "exclude", "filter nodes that exclude")

	flag.Int64Var(&config.Retry, "retry", 3, "set speedtest retry")
	flag.Int64Var(&config.Timeout, "timeout", 20, "set speedtest timeout")
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

	if config.Timeout < 10 {
		config.Timeout = 10
	}
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
		if len(config.Include) > 0 && !util.StringContains(proxy.Name, []string(config.Include)...) {
			continue
		}
		if len(config.Exclude) > 0 && util.StringContains(proxy.Name, []string(config.Exclude)...) {
			continue
		}
		names = append(names, proxy.Name)
	}
	sort.Strings(names)

	if len(names) == 0 {
		logx.Fatal("[config] no nodes left, please change include and exclude arguments")
	}

	logx.Infof("[speedtest] total nodes: %d", len(names))
	for i := 0; i < len(names); i++ {
		logx.Infof("-> %s", names[i])
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

		result, node := &Result{}, &Node{Name: proxy.Name}

		for j := int64(1); j <= config.Retry; j++ {
			data, err := fast.GetData()
			if err != nil {
				logx.WithField("err", err).Errorf("[fast.com] api fail")
				logx.Warnf("[speedtest] attempts %d time(s)", j)
				time.Sleep(time.Second)
				continue
			}

			logx.Infof("[%s] (%s) country: %s, city: %s", proxy.Name, data.Client.IP, data.Client.Location.Country, data.Client.Location.City)

			if len(data.Targets) == 0 {
				logx.Errorf("[%s] current area not exist speedtest node", proxy.Name)
				break
			}

			target := data.Targets[0]

			logx.Infof("[%s] speedtest node country: %s, city: %s", proxy.Name, target.Location.Country, target.Location.City)

			result, err = SpeedTest(target.URL, config.Process)
			if err != nil {
				logx.WithField("err", err).Errorf("[%s] speedtest fail", proxy.Name)
				continue
			}

			kb := float64(result.TotalBytes) / 1024
			ti := float64(result.TotalTime) / float64(time.Second)
			logx.Infof("[%s] speedtest download: %d kb, took: %.03f s, speed: %.02f kb/s", proxy.Name, int64(kb), ti, kb/ti)
			node.Speed = kb / ti
			break
		}

		util.ProxySet("", "")

		if result != nil {
			dashboard.TotalBytes += result.TotalBytes
			dashboard.TotalTime += result.TotalTime
		}
		dashboard.Nodes[i] = node

		logx.Infof("[%s] speedtest done, %d/%d", proxy.Name, i+1, len(names))
	}

	logx.Infof("total bytes: %.02f mb, total time: %d s", float64(dashboard.TotalBytes)/1024/1024, int64(dashboard.TotalTime/time.Second))

	bb := &bytes.Buffer{}
	writer := tablewriter.NewWriter(bb)
	writer.SetAutoFormatHeaders(false)
	writer.SetHeader([]string{"name", "speed(kb/s)"})
	writer.SetFooter([]string{"name", "speed(kb/s)"})
	writer.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT})
	for i := 0; i < len(names); i++ {
		node, color := dashboard.Nodes[i], tablewriter.FgGreenColor
		if node.Speed < 1024 {
			color = tablewriter.FgRedColor
		} else if node.Speed < 3072 {
			color = tablewriter.FgYellowColor
		}
		writer.Rich(
			[]string{node.Name, fmt.Sprintf("%.02f", node.Speed)},
			[]tablewriter.Colors{tablewriter.Color(tablewriter.Bold), tablewriter.Color(tablewriter.Bold, color)},
		)
	}
	writer.Render()

	ss := strings.Split(bb.String()[:bb.Len()-1], tablewriter.NEWLINE)
	for i := 0; i < len(ss); i++ {
		logx.Info(ss[i])
	}
}
