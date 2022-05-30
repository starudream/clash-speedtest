package clash

type Mode string

const (
	ModeGlobal Mode = "global"
	ModeRule   Mode = "rule"
	ModeDirect Mode = "direct"
)

var ModeMap = map[string]Mode{
	string(ModeGlobal): ModeGlobal,
	string(ModeRule):   ModeRule,
	string(ModeDirect): ModeDirect,
}
