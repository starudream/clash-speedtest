package clash

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestClient_GetVersion(t *testing.T) {
	version, err := c.GetVersion()
	testutil.LogNoErr(t, err, version)
}
