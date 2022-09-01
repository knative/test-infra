package shx

// CommandBuilder creates PreparedCommand's with common configuration
// such as always stopping on errors, running a set of commands in a
// directory, or using a set of environment variables.
type CommandBuilder struct {
	StopOnError bool
	Env         []string
	Dir         string
}

// Command creates a command using common configuration.
func (b *CommandBuilder) Command(cmd string, args ...string) PreparedCommand {
	return Command(cmd, args...).
		Must(b.StopOnError).
		Env(b.Env...).
		In(b.Dir)
}

// Run the given command, directing stderr to this program's stderr and
// printing stdout to stdout if mage was run with -v.
func (b *CommandBuilder) Run(cmd string, args ...string) error {
	return b.Command(cmd, args...).Run()
}

// RunS is like Run, but the command output is not written to stdout/stderr.
func (b *CommandBuilder) RunS(cmd string, args ...string) error {
	return b.Command(cmd, args...).RunS()
}

// RunE is like Run, but it only writes the command's output to os.Stderr when it fails.
func (b *CommandBuilder) RunE(cmd string, args ...string) error {
	return b.Command(cmd, args...).RunE()
}

// RunV is like Run, but always writes the command's stdout to os.Stdout.
func (b *CommandBuilder) RunV(cmd string, args ...string) error {
	return b.Command(cmd, args...).RunV()
}

// Output executes the prepared command, returning stdout.
func (b *CommandBuilder) Output(cmd string, args ...string) (string, error) {
	return b.Command(cmd, args...).Output()
}

// Outputs is like Output, but nothing is written to stdout/stderr.
func (b *CommandBuilder) OutputS(cmd string, args ...string) (string, error) {
	return b.Command(cmd, args...).OutputS()
}

// OutputE is like Output, but it only writes the command's output to os.Stderr when it fails.
func (b *CommandBuilder) OutputE(cmd string, args ...string) (string, error) {
	return b.Command(cmd, args...).OutputE()
}

// OutputV is like Output, but it always writes the command's stdout to os.Stdout.
func (b *CommandBuilder) OutputV(cmd string, args ...string) (string, error) {
	return b.Command(cmd, args...).OutputV()
}
