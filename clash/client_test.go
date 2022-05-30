package clash

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/starudream/clash-speedtest/internal/itest"
	"github.com/starudream/clash-speedtest/internal/json"
)

var x = New("http://127.0.0.1:9090", "")

func TestMain(m *testing.M) {
	itest.Init(m)
}

func TestClient_GetVersion(t *testing.T) {
	resp, err := x.GetVersion()
	require.NoError(t, err)
	t.Log(json.MustMarshalString(resp))
}

func TestClient_GetMode(t *testing.T) {
	resp, err := x.GetMode()
	require.NoError(t, err)
	t.Log(json.MustMarshalString(resp))
}

func TestClient_SetMode(t *testing.T) {
	err := x.SetMode(ModeRule)
	require.NoError(t, err)
}

func TestClient_GetProxies(t *testing.T) {
	resp, err := x.GetProxies()
	require.NoError(t, err)
	t.Log(json.MustMarshalString(resp))
}

func TestClient_SetProxy(t *testing.T) {
	err := x.SetProxy("🇭🇰 Premium|广港|IEPL|06")
	require.NoError(t, err)
}
