package clash

import (
	"github.com/starudream/go-lib/core/v2/config"
)

var c = NewClient(config.Get("clash.addr").String(), config.Get("clash.secret").String())
