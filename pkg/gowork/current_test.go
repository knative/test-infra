package gowork_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"knative.dev/test-infra/pkg/gowork"
	"knative.dev/test-infra/pkg/gowork/testdata"
)

func TestCurrent(t *testing.T) {
	t.Parallel()
	tcs := []currentTestCase{{
		name: "using GO111MODULE equals off",
		dir:  "code/foo",
		env:  testdata.Env{"GO111MODULE": "off"},
		err:  gowork.ErrInvalidGomod,
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
		err:  gowork.ErrInvalidGomod,
	}, {
		name: "in project without go.work",
		dir:  "code/bar",
		want: "bar",
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tfs := testdata.FS{Dir: tc.dir, Files: testdata.ExampleFS()}
			mod, err := gowork.Current(tfs, tc.env)
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
	env  testdata.Env
	err  error
	want string
}
