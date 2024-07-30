package emoji

import (
	"bytes"
	_ "embed"

	"github.com/rivo/uniseg"

	"github.com/starudream/go-lib/core/v2/codec/json"
	"github.com/starudream/go-lib/core/v2/utils/osutil"
)

type iItem struct {
	Emoji          string   `json:"emoji"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Aliases        []string `json:"aliases"`
	Tags           []string `json:"tags"`
	UnicodeVersion string   `json:"unicode_version"`
	IosVersion     string   `json:"ios_version"`
	SkinTones      bool     `json:"skin_tones,omitempty"`
}

var (
	//go:embed emoji.json
	emojiRaw []byte

	emojis = map[string]struct{}{}
)

func init() {
	items, err := json.UnmarshalTo[[]*iItem](emojiRaw)
	if err != nil {
		osutil.PanicErr(err)
	}
	for _, item := range items {
		emojis[item.Emoji] = struct{}{}
	}
}

func Remove(s string) string {
	bb := bytes.Buffer{}
	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		if _, ok := emojis[gr.Str()]; !ok {
			bb.WriteString(gr.Str())
			continue
		}
	}
	return bb.String()
}
