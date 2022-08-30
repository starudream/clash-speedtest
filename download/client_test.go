package download

import (
	"testing"
	"time"

	"github.com/starudream/go-lib/testx"
)

var x = New("http://127.0.0.1:7890", 3*time.Second)

func TestClient_Cloudflare(t *testing.T) {
	result, err := x.Cloudflare()
	testx.P(t, err, result)
}
