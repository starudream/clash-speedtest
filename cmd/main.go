package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/starudream/go-lib/app"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/clash-speedtest/clash"
	"github.com/starudream/clash-speedtest/download"
)

func init() {
	log.Attach("app", "clash-speedtest")

	config.SetDefault("url", "http://127.0.0.1:9090")
	config.SetDefault("proxy", "http://127.0.0.1:7890")
	config.SetDefault("retry", 3)
	config.SetDefault("timeout", "5s")
}

func main() {
	app.Add(initSpeedtest)
	app.Defer(deferSpeedtest)
	err := app.OnceGo()
	if err != nil {
		log.Fatal().Msgf("app init fail: %v", err)
	}
}

var (
	smu     sync.Mutex
	mode    clash.Mode
	results []*download.Result
)

func initSpeedtest(context.Context) error {
	clashCli := clash.New(config.GetString("url"), config.GetString("secret"))
	downloadCli := download.New(config.GetString("proxy"), config.GetDuration("timeout"))

	version, err := clashCli.GetVersion()
	if err != nil {
		return err
	}

	log.Info().Msgf("clash version: %s, premium: %t", version.Version, version.Premium)

	mode, err = clashCli.GetMode()
	if err != nil {
		return err
	}

	proxies, err := clashCli.GetProxies()
	if err != nil {
		return err
	}

	global := proxies["GLOBAL"]
	if global == nil || len(global.All) == 0 {
		return fmt.Errorf("clash no global proxy")
	}

	err = clashCli.SetMode(clash.ModeGlobal)
	if err != nil {
		return err
	}

	var (
		include = config.GetStringSlice("include")
		exclude = config.GetStringSlice("exclude")
		retry   = config.GetInt("retry")
	)

	if retry <= 0 {
		retry = 1
	}

	log.Info().Msgf("include: %#v, exclude: %#v", include, exclude)

	var names []string

	for i := 0; i < len(global.All); i++ {
		proxy, exist := proxies[global.All[i]]
		if !exist {
			return fmt.Errorf("clash proxy %s not exist", global.All[i])
		}

		l := log.With().Str("name", proxy.Name).Str("type", proxy.Type).Logger()

		if len(exclude) > 0 && includeKeywords(proxy.Name, exclude) {
			l.Debug().Msgf("exclude, skip proxy")
			continue
		}

		if len(include) > 0 && !includeKeywords(proxy.Name, include) {
			l.Debug().Msgf("include, skip proxy")
			continue
		}

		switch proxy.Type {
		case "Shadowsocks", "ShadowsocksR", "Snell", "Socks5", "Http", "Vmess", "Trojan":
		case "Direct", "Reject", "Relay", "Selector", "Fallback", "URLTest", "LoadBalance":
			continue
		default:
			l.Warn().Msgf("clash type %s not support", proxy.Type)
			continue
		}

		names = append(names, proxy.Name)

		l.Debug().Msgf("add proxy")
	}

	log.Info().Msgf("speedtest count: %d", len(names))

	if len(names) == 0 {
		log.Warn().Msgf("no proxy to test")
		return nil
	}

	defer deferSpeedtest()

	for i := 0; i < len(names); i++ {
		proxy := proxies[names[i]]

		l := log.With().Int("index", i+1).Str("name", proxy.Name).Str("type", proxy.Type).Logger()

		err = clashCli.SetProxy(proxy.Name)
		if err != nil {
			return err
		}

		l.Info().Msgf("speedtest begin")

		for j := 0; j < retry; j++ {
			result, se := downloadCli.Cloudflare()
			if se != nil {
				l.Error().Msgf("speedtest error: %v", se)
				continue
			}

			result.Name = proxy.Name
			result.Type = proxy.Type

			smu.Lock()
			results = append(results, result)
			smu.Unlock()

			l.Info().
				IPAddr("ip", net.ParseIP(result.IP)).
				Str("country", result.Country).Str("city", result.City).
				Str("lat", result.Lat).Str("lng", result.Lng).
				Msgf(
					"speedtest end, took %.02fs, duration %dms, download %s, speed %s",
					result.TotalDuration.Seconds(), result.ConnectDuration.Milliseconds(), formatMeter(float64(result.Length)), formatMeter(result.BS)+"/s",
				)

			break
		}
	}

	return nil
}

func deferSpeedtest() {
	func() {
		smu.Lock()
		defer smu.Unlock()
		if len(results) > 0 {
			text := formatResult(results)
			fmt.Println(text)
			filename := filepath.Join(config.GetString("output"), fmt.Sprintf("result-%s.txt", time.Now().Format("20060102150405")))
			err := os.WriteFile(filename, []byte(text), 0644)
			if err != nil {
				log.Fatal().Msgf("write file fail: %v", err)
			}
			clashCli := clash.New(config.GetString("url"), config.GetString("secret"))
			optimumNodeResult := results[0]
			err = clashCli.SetProxy(optimumNodeResult.Name)
			if err != nil {
				log.Info().Msgf("Selected node: %s for online serfing🏄🏻, bandwidth: %s, latency: %dms, enjoy!", optimumNodeResult.Name, formatMeter(optimumNodeResult.BS), optimumNodeResult.ConnectDuration.Milliseconds())
			} else {
				log.Error().Msgf("Select node: %s error, reason: %s", optimumNodeResult.Name, err)
			}
			results = nil
		}
	}()
	func() {
		if mode != "" {
			_ = clash.New(config.GetString("url"), config.GetString("secret")).SetMode(mode)
			mode = ""
		}
	}()
}

func formatResult(results []*download.Result) string {
	sort.Slice(results, func(lhs, rhs int) bool {
		return results[lhs].BS >= results[rhs].BS
	})
	bb := &bytes.Buffer{}
	tw := tablewriter.NewWriter(bb)
	tw.SetAlignment(tablewriter.ALIGN_CENTER)
	tw.SetHeader([]string{"name", "country", "type", "duration", "speed", "total"})

	total := float64(0)
	for i := 0; i < len(results); i++ {
		v := results[i]
		total += float64(v.Length)
		var color int
		if v.BS < 1024*1024 {
			color = tablewriter.FgRedColor
		} else if v.BS < 5*1024*1024 {
			color = tablewriter.FgYellowColor
		} else {
			color = tablewriter.FgGreenColor
		}
		tw.Rich(
			[]string{v.Name, v.Country, v.Type, fmt.Sprintf("%dms", v.ConnectDuration.Milliseconds()), formatMeter(v.BS) + "/s", formatMeter(float64(v.Length))},
			[]tablewriter.Colors{{tablewriter.Bold}, {}, {}, {}, {tablewriter.Bold, color}, {}},
		)
	}

	tw.SetFooter([]string{"", "", "", "", "", formatMeter(total)})
	tw.Render()
	return bb.String()
}

func formatMeter(v float64) string {
	if v < 1024 {
		return fmt.Sprintf("%.02fB", v)
	}
	v /= 1024
	if v < 1024 {
		return fmt.Sprintf("%.02fKB", v)
	}
	v /= 1024
	if v < 1024 {
		return fmt.Sprintf("%.02fMB", v)
	}
	v /= 1024
	if v < 1024 {
		return fmt.Sprintf("%.02fGB", v)
	}
	v /= 1024
	return fmt.Sprintf("%.02fTB", v)
}

func includeKeywords(v string, keywords []string) bool {
	for _, k := range keywords {
		if strings.Contains(v, k) {
			return true
		}
	}
	return false
}
