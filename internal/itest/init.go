package itest

import (
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/starudream/clash-speedtest/internal/ilog"
)

func Init(m *testing.M) {
	log.Logger = log.Output(ilog.NewConsoleWriter())

	viper.SetEnvPrefix("scs")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("debug", true)

	os.Exit(m.Run())
}
