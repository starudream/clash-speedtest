package clash

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestClient_GetProxies(t *testing.T) {
	proxies1, err := c.GetProxies()
	testutil.Nil(t, err)
	proxies2 := proxies1.FilterProxies()
	testutil.Log(t, len(proxies2), proxies2)
}

func TestClient_SetGlobalProxy(t *testing.T) {
	testutil.Nil(t, c.SetGlobalProxy(ProxyDirect))
}
