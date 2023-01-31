package cli_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"knative.dev/test-infra/pkg/gowork/testdata"
	"knative.dev/test-infra/tools/modscope/cli"
)

func TestListCmd(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	app := cli.App{}
	c := app.Command()
	c.SetOut(&buf)
	c.SetArgs([]string{"list"})
	err := c.Execute()

	assert.NoError(t, err)
	assert.Equal(t, "knative.dev/test-infra\n", buf.String())
}

func TestList(t *testing.T) {
	t.Parallel()
	tcs := []testListCase{{
		name: "in project's root dir",
		dir:  "code/foo",
		want: []string{"foo", "foo/a"},
	}, {
		name: "in project's root dir showing paths",
		dir:  "code/foo",
		fl:   cli.Flags{DisplayFilepath: true},
		want: []string{"/code/foo", "/code/foo/a"},
	}, {
		name: "in non go.work project",
		dir:  "code/bar",
		want: []string{"bar"},
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			os := testOs{
				FS: testdata.FS{
					Dir:   tc.dir,
					Files: testdata.ExampleFS(),
				},
			}
			printer := &testPrinter{}
			err := cli.List(os, &tc.fl, printer)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				want := strings.Join(tc.want, "\n") + "\n"
				assert.Equal(t, want, printer.buf.String())
			}
		})
	}
}

type testListCase struct {
	name string
	fl   cli.Flags
	dir  string
	err  error
	want []string
}

type testOs struct {
	testdata.FS
	testdata.Env
}

func (t testOs) Abs(filepath string) string {
	return "/" + filepath
}

type testPrinter struct {
	buf bytes.Buffer
}

func (t *testPrinter) Println(i ...interface{}) {
	_, err := fmt.Fprintln(&t.buf, i...)
	if err != nil {
		panic(err)
	}
}
