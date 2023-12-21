package speedtest

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestClient_Servers(t *testing.T) {
	servers, err := c.GetServers()
	testutil.LogNoErr(t, err, servers)

	for i := 0; i < min(3, len(servers)); i++ {
		latency, err := c.LatencyServer(servers[i])
		testutil.LogNoErr(t, err, latency.String())
	}
}
