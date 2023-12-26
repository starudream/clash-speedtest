package cloudflare

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestClient_Download(t *testing.T) {
	res, err := c.Download(0, nil)
	testutil.LogNoErr(t, err, res.String())
}
