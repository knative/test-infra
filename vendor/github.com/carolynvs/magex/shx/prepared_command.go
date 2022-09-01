package shx

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/carolynvs/magex/mgx"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type PreparedCommand struct {
	Cmd         *exec.Cmd
	StopOnError bool
}

// Command creates a default command. Stdout is logged in verbose mode. Stderr
// is sent to os.Stderr.
func Command(cmd string, args ...string) PreparedCommand {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = os.Environ()
	return PreparedCommand{Cmd: c}
}

// String prints the command-line representation of the PreparedCommand.
func (c PreparedCommand) String() string {
	return strings.Join(c.Cmd.Args, " ")
}

// Must immediately stops the build when the command fails.
func (c PreparedCommand) Must(stopOnError ...bool) PreparedCommand {
	switch len(stopOnError) {
	case 0:
		c.StopOnError = true
	case 1:
		c.StopOnError = stopOnError[0]
	default:
		mgx.Must(fmt.Errorf("More than one value for Must(stopOnError ...string) was passed to the command %s", c))
	}
	return c
}

// Args appends additional arguments to the command.
func (c PreparedCommand) Args(args ...string) PreparedCommand {
	c.Cmd.Args = append(c.Cmd.Args, args...)
	return c
}

// CollapseArgs removes empty arguments from the argument list.
//
// This is helpful when sometimes a flag should be specified and
// sometimes it shouldn't.
func (c PreparedCommand) CollapseArgs() PreparedCommand {
	result := make([]string, 0, len(c.Cmd.Args))
	for _, arg := range c.Cmd.Args {
		if arg != "" {
			result = append(result, arg)
		}
	}
	c.Cmd.Args = result
	return c
}

// Env defines additional environment variables for the command.
// All ambient environment variables are included by default.
// Example:
//  c.Env("X=1", "Y=2")
func (c PreparedCommand) Env(vars ...string) PreparedCommand {
	for _, v := range vars {
		c.Cmd.Env = append(c.Cmd.Env, v)
	}
	return c
}

// In sets the working directory of the command.
func (c PreparedCommand) In(dir string) PreparedCommand {
	c.Cmd.Dir = dir
	return c
}

// Stdin sets the command's stdin.
func (c PreparedCommand) Stdin(stdin io.Reader) PreparedCommand {
	c.Cmd.Stdin = stdin
	return c
}

// Stdout directs stdout from the command.
func (c PreparedCommand) Stdout(stdout io.Writer) PreparedCommand {
	c.Cmd.Stdout = stdout
	return c
}

// Stderr directs stderr from the command.
func (c PreparedCommand) Stderr(stderr io.Writer) PreparedCommand {
	c.Cmd.Stderr = stderr
	return c
}

// Runs a command silently, without writing to stdout/stderr.
func (c PreparedCommand) Silent() PreparedCommand {
	c.Cmd.Stdout = nil
	c.Cmd.Stderr = nil
	return c
}

// Exec the prepared command, returning if the command was run and its
// exit code. Does not modify the configured outputs.
func (c PreparedCommand) Exec() (ran bool, code int, err error) {
	if mg.Verbose() {
		log.Println("exec:", c.Cmd.Path, strings.Join(c.Cmd.Args, " "))
	}

	err = c.Cmd.Run()
	ran = sh.CmdRan(err)
	code = sh.ExitStatus(err)

	if err != nil {
		if ran {
			err = mg.Fatalf(code, `running "%s" failed with exit code %d`, c, code)
		} else {
			err = fmt.Errorf(`failed to run "%s: %v"`, c, err)
		}
		if c.StopOnError {
			mgx.Must(err)
		}
	}

	return ran, code, err
}

// Run the given command, directing stderr to os.Stderr and
// printing stdout to os.Stdout if mage was run with -v.
func (c PreparedCommand) Run() error {
	if mg.Verbose() {
		c.Cmd.Stdout = os.Stdout
	} else {
		c.Cmd.Stdout = nil
	}

	_, _, err := c.Exec()
	return err
}

// RunV is like Run, but always writes the command output to os.Stdout.
func (c PreparedCommand) RunV() error {
	c.Stdout(os.Stdout)
	_, _, err := c.Exec()
	return err
}

// RunE is like Run, but it only writes the command combined to os.Stderr when it fails.
func (c PreparedCommand) RunE() error {
	output := &bytes.Buffer{}
	c.Stdout(output)
	c.Stderr(output)
	_, _, err := c.Exec()
	if err != nil {
		fmt.Fprint(os.Stderr, output.String())
	}
	return err
}

// RunS is like Run, but the command output is not written to stdout/stderr.
func (c PreparedCommand) RunS() error {
	_, _, err := c.Silent().Exec()
	return err
}

// Output executes the prepared command, directing stderr to os.Stderr and
// printing stdout to os.Stdout if mage was run with -v. The command's stdout
// is always returned.
func (c PreparedCommand) Output() (string, error) {
	stdout := &bytes.Buffer{}
	if mg.Verbose() {
		c.Cmd.Stdout = io.MultiWriter(stdout, os.Stdout)
	} else {
		c.Cmd.Stdout = stdout
	}

	_, _, err := c.Exec()
	return strings.TrimSuffix(stdout.String(), "\n"), err
}

// OutputV is like Output, but it always writes the command output to os.Stdout.
func (c PreparedCommand) OutputV() (string, error) {
	stdout := &bytes.Buffer{}
	c.Cmd.Stdout = io.MultiWriter(stdout, os.Stdout)
	_, _, err := c.Exec()
	return strings.TrimSuffix(stdout.String(), "\n"), err
}

// OutputE is like Output, but it only writes the command output to os.Stderr when it fails.
func (c PreparedCommand) OutputE() (string, error) {
	stdout := &bytes.Buffer{}
	output := &bytes.Buffer{}
	c.Stdout(io.MultiWriter(stdout, output))
	c.Stderr(output)
	_, _, err := c.Exec()
	if err != nil {
		fmt.Fprint(os.Stderr, output.String())
	}
	return strings.TrimSuffix(stdout.String(), "\n"), err
}

// Outputs is like Output, but the command output is not written to stdout/stderr.
func (c PreparedCommand) OutputS() (string, error) {
	stdout := &bytes.Buffer{}
	_, _, err := c.Stdout(stdout).Stderr(nil).Exec()
	return strings.TrimSuffix(stdout.String(), "\n"), err
}
