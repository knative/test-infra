package shx

var defaultCommand CommandBuilder

// Run the given command, directing stderr to this program's stderr and
// printing stdout to stdout if mage was run with -v.
func Run(cmd string, args ...string) error {
	return defaultCommand.Run(cmd, args...)
}

// RunS is like Run, but the command output is not written to stdout/stderr.
func RunS(cmd string, args ...string) error {
	return defaultCommand.RunS(cmd, args...)
}

// RunE is like Run, but it only writes the command's output to os.Stderr when it fails.
func RunE(cmd string, args ...string) error {
	return defaultCommand.RunE(cmd, args...)
}

// RunV is like Run, but always writes the command's stdout to os.Stdout.
func RunV(cmd string, args ...string) error {
	return defaultCommand.RunV(cmd, args...)
}

// Output executes the prepared command, returning stdout.
func Output(cmd string, args ...string) (string, error) {
	return defaultCommand.Output(cmd, args...)
}

// Outputs is like Output, but nothing is written to stdout/stderr.
func OutputS(cmd string, args ...string) (string, error) {
	return defaultCommand.OutputS(cmd, args...)
}

// OutputE is like Output, but it only writes the command's output to os.Stderr when it fails.
func OutputE(cmd string, args ...string) (string, error) {
	return defaultCommand.OutputE(cmd, args...)
}

// OutputV is like Output, but it always writes the command's stdout to os.Stdout.
func OutputV(cmd string, args ...string) (string, error) {
	return defaultCommand.OutputV(cmd, args...)
}
