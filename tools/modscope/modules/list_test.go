package modules_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"knative.dev/test-infra/tools/modscope/modules"
	"knative.dev/test-infra/tools/modscope/test"
)

func TestList(t *testing.T) {
	want := []string{"foo", "foo/a"}
	t.Parallel()
	tcs := []listTestCase{{
		name: "using GOWORK equals off",
		dir:  "code/foo",
		env:  test.Env{"GOWORK": "off"},
		err:  modules.ErrInvalidGowork,
	}, {
		name: "using GOWORK to point to go.work from outside of project dir",
		dir:  "srv",
		env:  test.Env{"GOWORK": "/code/foo/go.work"},
	}, {
		name: "in project's root dir",
		dir:  "code/foo",
	}, {
		name: "in project subdir",
		dir:  "code/foo/a",
	}, {
		name: "in project subdir not listed in go.work",
		dir:  "code/foo/b",
	}, {
		name: "outside of project dir",
		err:  modules.ErrInvalidGowork,
	}, {
		name: "in project without go.work",
		dir:  "code/bar",
		err:  modules.ErrInvalidGowork,
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tfs := test.FS{Dir: tc.dir, Files: test.ExampleFS()}
			mods, err := modules.List(tfs, tc.env)
			got := toNames(mods)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

type listTestCase struct {
	name string
	dir  string
	env  test.Env
	err  error
}

func toNames(mods []modules.Module) []string {
	names := make([]string, len(mods))
	for i, mod := range mods {
		names[i] = mod.Name
	}
	sort.Strings(names)
	return names
}
