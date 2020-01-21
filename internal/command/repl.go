package command

import (
	"os"

	"github.com/axetroy/dvs/internal/vm"
)

type ReplOption struct {
	Image string
}

func Repl(option *ReplOption) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	v, err := vm.NewVirtualMachine(&vm.Options{
		Cwd:   cwd,
		Image: option.Image,
	})

	if err != nil {
		return err
	}

	if err := v.Start(); err != nil {
		return err
	}

	defer v.Destroy()

	if err := v.Attach(); err != nil {
		return err
	}

	return nil
}
