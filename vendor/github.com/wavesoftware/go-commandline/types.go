package commandline

import (
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/wavesoftware/go-retcode"
)

// ErrNoRootCommand is returned when no root command is provided.
var ErrNoRootCommand = errors.New("no root command provided")

// App represents a command line application.
type App struct {
	CobraProvider
	Exit func(code int)
	root *cobra.Command
}

// CobraProvider is used to provide a Cobra command.
type CobraProvider interface {
	Command() *cobra.Command
}

// Option is used to configure an App.
type Option func(*App)

// New creates a new App from CobraProvider.
func New(cp CobraProvider) *App {
	return &App{
		CobraProvider: cp,
		Exit:          os.Exit,
	}
}

// ExecuteOrDie will execute the application or perform os.Exit in case of error.
func (a *App) ExecuteOrDie(options ...Option) {
	if err := a.Execute(options...); err != nil {
		a.Exit(retcode.Calc(err))
	}
}

// WithArgs creates an option which sets args.
func WithArgs(args ...string) Option {
	return func(app *App) {
		app.root.SetArgs(args)
	}
}

// WithInput creates an option witch sets os.Stdin.
func WithInput(in io.Reader) Option {
	return func(app *App) {
		app.root.SetIn(in)
	}
}

// WithOutput creates an option witch sets os.Stdout and os.Stderr.
func WithOutput(out io.Writer) Option {
	return func(app *App) {
		app.root.SetOut(out)
		app.root.SetErr(out)
	}
}

// WithExit creates an option which sets the exit function.
func WithExit(fn func(code int)) Option {
	return func(app *App) {
		app.Exit = fn
	}
}

// Execute will execute the application with the provided options and return
// error if any.
func (a *App) Execute(options ...Option) error {
	if err := a.init(); err != nil {
		return err
	}
	for _, config := range options {
		config(a)
	}
	// cobra.Command should pass our own errors, no need to wrap them.
	return a.root.Execute() //nolint:wrapcheck
}

func (a *App) init() error {
	if a.Exit == nil {
		a.Exit = os.Exit
	}
	if a.CobraProvider == nil {
		return ErrNoRootCommand
	}
	a.root = a.Command()
	if a.root == nil {
		return ErrNoRootCommand
	}
	return nil
}
