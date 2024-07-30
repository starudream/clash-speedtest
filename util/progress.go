package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"

	"github.com/starudream/go-lib/core/v2/slog"

	"github.com/starudream/clash-speedtest/emoji"
)

type ProgressBar = pb.ProgressBar

const tmpl = `{{with string . "prefix"}}{{.}} {{end}}{{counters . }} {{bar . "[" "=" ">" "-" "]"}} {{percent . }}`

func NewBarsPool(n int, name string) (*pb.Pool, []*pb.ProgressBar) {
	bars := make([]*pb.ProgressBar, n)
	for i := 0; i < n; i++ {
		bars[i] = pb.New64(0).
			SetRefreshRate(50*time.Millisecond).
			SetTemplateString(tmpl).
			Set(pb.SIBytesPrefix, true).
			Set("prefix", fmt.Sprintf("%s [%d]", strings.TrimSpace(emoji.Remove(name)), i+1))
	}
	pool, err := pb.StartPool(bars...)
	if err != nil {
		slog.Error("start progress bars error: %v", err)
	}
	return pool, bars
}
