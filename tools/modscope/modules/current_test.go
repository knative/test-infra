package modules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"knative.dev/test-infra/tools/modscope/modules"
	"knative.dev/test-infra/tools/modscope/test"
)

func TestCurrent(t *testing.T) {
	t.Parallel()
	tcs := []currentTestCase{{
		name: "using GO111MODULE equals off",
		dir:  "code/foo",
		env:  test.Env{"GO111MODULE": "off"},
		err:  modules.ErrInvalidGomod,
	}, {
		name: "in project's root dir",
		dir:  "code/foo",
		want: "foo",
	}, {
		name: "in project subdir",
		dir:  "code/foo/a",
		want: "foo/a",
	}, {
		name: "outside of project dir",
		err:  modules.ErrInvalidGomod,
	}, {
		name: "in project without go.work",
		dir:  "code/bar",
		want: "bar",
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tfs := test.FS{Dir: tc.dir, Files: test.ExampleFS()}
			mod, err := modules.Current(tfs, tc.env)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, mod.Name)
			}
		})
	}
}

type currentTestCase struct {
	name string
	dir  string
	env  test.Env
	err  error
	want string
}
