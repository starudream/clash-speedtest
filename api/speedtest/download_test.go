package speedtest

import (
	"math/rand"
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestClient_Download(t *testing.T) {
	servers, err := c.GetServers()
	testutil.LogNoErr(t, err, servers)
	testutil.NotEqual(t, 0, len(servers))

	server := servers[rand.Intn(len(servers))]
	testutil.Log(t, server)

	res, err := c.Download(server, 0, nil)
	testutil.LogNoErr(t, err, res.String())
}
