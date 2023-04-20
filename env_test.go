package cleo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Env_Get_Default(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		env  *Env
	}{
		{"nil env", nil},
		{"nil map", &Env{}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			act := tc.env.Get("PATH")
			r.True(len(act) > 0)
		})
	}

}

func Test_Env_Get_Map(t *testing.T) {
	t.Parallel()

	empty := &Env{
		data: map[string]string{},
	}

	bin := "/usr/local/bin"
	full := &Env{
		data: map[string]string{
			"PATH": bin,
		},
	}

	tcs := []struct {
		name string
		env  *Env
		exp  string
	}{
		{"empty map", empty, ""},
		{"full map", full, bin},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			act := tc.env.Get("PATH")
			r.Equal(tc.exp, act)
		})
	}
}
