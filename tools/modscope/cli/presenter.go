package cli

import (
	"knative.dev/test-infra/tools/modscope/modules"
)

// Printer is an interface for printing the output.
type Printer interface {
	Println(...interface{})
}

type presenter struct {
	OS
	*Flags
	Printer
}

func (p presenter) presentList(mods []modules.Module, err error) error {
	if err != nil {
		return err
	}
	for _, m := range mods {
		_ = p.presentModule(&m, nil)
	}
	return nil
}

func (p presenter) presentModule(curr *modules.Module, err error) error {
	if err != nil {
		return err
	}
	if curr == nil {
		return nil
	}
	line := curr.Name
	if p.DisplayFilepath {
		line = p.Abs(curr.Path)
	}
	p.Println(line)
	return nil
}
