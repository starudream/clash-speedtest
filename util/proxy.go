package util

import (
	"os"
)

const (
	httpProxyName  = "HTTP_PROXY"
	httpsProxyName = "HTTPS_PROXY"
)

func ProxySet(hp, hsp string) {
	_ = os.Setenv(httpProxyName, hp)
	_ = os.Setenv(httpsProxyName, hsp)
}

func ProxyGet() (string, string) {
	return os.Getenv(httpProxyName), os.Getenv(httpsProxyName)
}
