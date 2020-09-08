package kubetest2

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func run(command *exec.Cmd) (string, error) {
	var buf bytes.Buffer
	command.Stdout = io.MultiWriter(&buf, os.Stdout)
	command.Stderr = io.MultiWriter(&buf, os.Stderr)
	err := command.Run()
	return buf.String(), err
}
