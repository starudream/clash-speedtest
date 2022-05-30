package json

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMustUnmarshalTo(t *testing.T) {
	type X struct {
		A string
		B int
		C float64
		D bool
		E time.Time
		F []byte
	}

	var x = &X{
		A: "foo",
		B: 5,
		C: 3.14,
		D: true,
		E: time.Now().Truncate(time.Second),
		F: []byte("bar"),
	}

	t.Logf("%#v", x)

	bs := MustMarshal(x)

	t.Log(string(bs))

	xx, err := UnmarshalTo[*X](bs)
	require.NoError(t, err)

	t.Logf("%#v", xx)

	xy := MustUnmarshalTo[*X](bs)

	t.Logf("%#v", xy)

	require.Equal(t, x.A, xy.A)
	require.Equal(t, x.B, xy.B)
	require.Equal(t, x.C, xy.C)
	require.Equal(t, x.D, xy.D)
	require.Equal(t, x.E, xy.E)
	require.Equal(t, x.F, xy.F)
}
