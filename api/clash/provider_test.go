package clash

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestClient_GetProviderProxies(t *testing.T) {
	providers, err := c.GetProviderProxies()
	testutil.Nil(t, err)
	proxies := providers.FilterProxies()
	testutil.Log(t, len(proxies), proxies)
}
