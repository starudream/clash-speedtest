package fast

import (
	"math/rand"
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestClient_Download(t *testing.T) {
	config, err := c.GetConfig()
	testutil.LogNoErr(t, err, config)
	testutil.NotEqual(t, 0, len(config.Targets))

	target := config.Targets[rand.Intn(len(config.Targets))]
	testutil.Log(t, target)

	res, err := c.Download(target, 0, nil)
	testutil.LogNoErr(t, err, res.String())
}
