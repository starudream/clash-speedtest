package speedtest

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestClient_GetConfig(t *testing.T) {
	config, err := c.GetConfig()
	testutil.LogNoErr(t, err, config)
}
