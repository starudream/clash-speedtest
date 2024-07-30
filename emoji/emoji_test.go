package emoji_test

import (
	"testing"

	"github.com/starudream/clash-speedtest/emoji"
)

func TestRemove(t *testing.T) {
	t.Log(emoji.Remove(`emoji 😀 😃`))
}
