package clash

import (
	"testing"

	"github.com/starudream/go-lib/testx"
)

var x = New("http://127.0.0.1:54615", "865aca09-4f3f-42b6-bc6e-fc270c54e75a")

func TestClient_GetVersion(t *testing.T) {
	resp, err := x.GetVersion()
	testx.P(t, err, resp)
}

func TestClient_GetMode(t *testing.T) {
	resp, err := x.GetMode()
	testx.P(t, err, resp)
}

func TestClient_SetMode(t *testing.T) {
	err := x.SetMode(ModeRule)
	testx.P(t, err)
}

func TestClient_GetProxies(t *testing.T) {
	resp, err := x.GetProxies()
	testx.P(t, err, resp)
}

func TestClient_SetProxy(t *testing.T) {
	err := x.SetProxy("🇭🇰 Premium|广港|IEPL|06")
	testx.P(t, err)
}
