package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/starudream/clash-speedtest/clash"
	"github.com/starudream/clash-speedtest/config"
	"github.com/starudream/clash-speedtest/download"
	"github.com/starudream/clash-speedtest/internal/app"
	"github.com/starudream/clash-speedtest/internal/ilog"
)

var rootCmd = &cobra.Command{
	Use:     config.AppName,
	Short:   config.AppName,
	Version: config.FULL_VERSION,
	Run: func(cmd *cobra.Command, args []string) {
		app.Add(initSpeedtest)
		app.Recover(deferSpeedtest)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	smu     sync.Mutex
	mode    clash.Mode
	results []*download.Result
)

func initSpeedtest(context.Context) error {
	clashCli := clash.New(viper.GetString("url"), viper.GetString("secret"))
	downloadCli := download.New(viper.GetString("proxy"), viper.GetDuration("timeout"))

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
		include = viper.GetStringSlice("include")
		exclude = viper.GetStringSlice("exclude")
		retry   = viper.GetInt("retry")
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
		case "Shadowsocks", "ShadowsocksR":
		case "Direct", "Reject", "Selector", "URLTest":
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
			filename := filepath.Join(viper.GetString("output"), fmt.Sprintf("result-%s.txt", time.Now().Format("20060102150405")))
			ilog.WrapError(os.WriteFile(filename, []byte(text), 0644))
			results = nil
		}
	}()
	func() {
		if mode != "" {
			ilog.WrapError(clash.New(viper.GetString("url"), viper.GetString("secret")).SetMode(mode))
			mode = ""
		}
	}()
}

func formatResult(results []*download.Result) string {
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
