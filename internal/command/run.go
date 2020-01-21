package command

import (
	"os"

	"github.com/axetroy/dvs/internal/vm"
)

type RunOption struct {
	Image string
}

func Run(command []string, option *RunOption) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	v, err := vm.NewVirtualMachine(&vm.Options{
		Cwd:      cwd,
		Image:    option.Image,
		Commands: &command,
	})

	if err != nil {
		return err
	}

	if err := v.Start(); err != nil {
		return err
	}

	defer v.Destroy()

	if err := v.Wait(); err != nil {
		return err
	}

	if err := v.Log(); err != nil {
		return err
	}

	return nil
}
