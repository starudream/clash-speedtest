package main

import (
	"github.com/starudream/clash-speedtest/cmd"
	"github.com/starudream/clash-speedtest/internal/app"
	"github.com/starudream/clash-speedtest/internal/ierr"
)

func main() {
	defer ierr.Recover(app.Stop)

	cmd.Execute()

	ierr.CheckErr(app.OnceGo())
}
