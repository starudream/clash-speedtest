package ios

import (
	"os"

	"github.com/starudream/clash-speedtest/internal/ierr"
)

var (
	exec string
	home string
	pwd  string
)

func init() {
	var err error

	exec, err = os.Executable()
	ierr.CheckErr(err)

	home, err = os.UserHomeDir()
	ierr.CheckErr(err)

	pwd, err = os.Getwd()
	ierr.CheckErr(err)
}

func Executable() string {
	return exec
}

func UserHomeDir() string {
	return home
}

func PWD() string {
	return pwd
}
