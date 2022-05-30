package ierr

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/viper"
)

func CheckErr(msg any) {
	if msg != nil {
		fmt.Fprintln(os.Stderr, "Error:", msg)
		if viper.GetBool("debug") {
			fmt.Fprintln(os.Stderr, "Stack:", string(debug.Stack()))
		}
		os.Exit(1)
	}
}
