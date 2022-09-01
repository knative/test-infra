package shx

import (
	"bytes"
	"io"
	"os"
)

type capture struct {
	f        *os.File
	releaseF os.File
	w        *os.File
	out      chan string
}

// RecordStdout records what is written to os.Stdout.
func RecordStdout() *capture {
	return captureFile(os.Stdout)
}

// RecordStderr records what is written to os.Stderr.
func RecordStderr() *capture {
	return captureFile(os.Stderr)
}

func captureFile(f *os.File) *capture {
	c := &capture{
		f:        f,
		releaseF: *f,
	}
	r, w, _ := os.Pipe()
	c.w = w
	*f = *w

	c.out = make(chan string)
	go func() {
		var buf bytes.Buffer
		io.MultiWriter(f, &buf)
		io.Copy(&buf, r)
		c.out <- buf.String()
	}()

	return c
}

// Release reverts the changes to the captured output.
func (c *capture) Release() {
	*c.f = c.releaseF
	c.w.Close()
}

// Output releases the captured file and returns the output.
func (c *capture) Output() string {
	c.Release()
	return <-c.out
}
