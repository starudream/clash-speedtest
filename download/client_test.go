package download

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/starudream/clash-speedtest/internal/itest"
	"github.com/starudream/clash-speedtest/internal/json"
)

var x = New("http://127.0.0.1:7890", 3*time.Second)

func TestMain(m *testing.M) {
	itest.Init(m)
}

func TestClient_Cloudflare(t *testing.T) {
	result, err := x.Cloudflare()
	require.NoError(t, err)
	t.Log(json.MustMarshalString(result))
}
