package cmd

import (
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"github.com/starudream/clash-speedtest/internal/ierr"
	"github.com/starudream/clash-speedtest/internal/ilog"
)

func init() {
	cobra.OnInitialize(initLogger, initConfig)

	rootCmd.PersistentFlags().BoolP("debug", "", false, "(env: SCS_DEBUG) show debug information")
	ierr.CheckErr(viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")))

	rootCmd.PersistentFlags().StringP("url", "", "http://127.0.0.1:9090", "(env: SCS_URL) clash external controller url")
	ierr.CheckErr(viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url")))

	rootCmd.PersistentFlags().StringP("secret", "", "", "(env: SCS_SECRET) clash external controller secret")
	ierr.CheckErr(viper.BindPFlag("secret", rootCmd.PersistentFlags().Lookup("secret")))

	rootCmd.PersistentFlags().StringP("proxy", "", "http://127.0.0.1:7890", "(env: SCS_PROXY) clash http proxy url")
	ierr.CheckErr(viper.BindPFlag("proxy", rootCmd.PersistentFlags().Lookup("proxy")))

	rootCmd.PersistentFlags().StringSliceP("include", "", []string{}, "(env: SCS_INCLUDE) filter nodes by include")
	ierr.CheckErr(viper.BindPFlag("include", rootCmd.PersistentFlags().Lookup("include")))

	rootCmd.PersistentFlags().StringSliceP("exclude", "", []string{}, "(env: SCS_EXCLUDE) filter nodes by exclude")
	ierr.CheckErr(viper.BindPFlag("exclude", rootCmd.PersistentFlags().Lookup("exclude")))

	rootCmd.PersistentFlags().IntP("retry", "", 3, "(env: SCS_RETRY) retry times when failed")
	ierr.CheckErr(viper.BindPFlag("retry", rootCmd.PersistentFlags().Lookup("retry")))

	rootCmd.PersistentFlags().DurationP("timeout", "", 5*time.Second, "(env: SCS_TIMEOUT) timeout for http request")
	ierr.CheckErr(viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout")))
}

func initConfig() {
	viper.SetEnvPrefix("scs") // starudream - clash - speedtest
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("debug", false)
	viper.SetDefault("log.level", "INFO")

	level, err := zerolog.ParseLevel(strings.ToLower(viper.GetString("log.level")))
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}

	debug := viper.GetBool("debug")
	if debug {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	if level < zerolog.NoLevel {
		log.Logger = log.Output(zerolog.MultiLevelWriter(ilog.NewConsoleWriter()))
		if debug {
			log.Logger = log.With().Caller().Logger()
		}
	}

	zerolog.DefaultContextLogger = &log.Logger
}

func initLogger() {
	w := ilog.New(log.Output(ilog.NewConsoleWriter()), "cfg")
	jww.TRACE = w.WithLevel(zerolog.TraceLevel)
	jww.DEBUG = w.WithLevel(zerolog.DebugLevel)
	jww.INFO = w.WithLevel(zerolog.InfoLevel)
	jww.WARN = w.WithLevel(zerolog.WarnLevel)
	jww.ERROR = w.WithLevel(zerolog.ErrorLevel)
	jww.CRITICAL = w.WithLevel(zerolog.FatalLevel)
	jww.FATAL = w.WithLevel(zerolog.FatalLevel)
	jww.LOG = w.WithLevel(zerolog.TraceLevel)
}
