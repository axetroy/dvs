package vm

import (
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type VirtualMachine struct {
	id     string          // docker container id
	ctx    context.Context // context
	client *client.Client  // docker client
}

type Options struct {
	Cwd      string    // current working dir
	Image    string    // the image name you want to run
	Commands *[]string // the COMMAND for image
}

func isImageExist(cli *client.Client, ctx context.Context, image string) (bool, error) {
	images, err := cli.ImageList(ctx, types.ImageListOptions{
		All: true,
	})

	if err != nil {
		return false, err
	}

	hasImageExist := false

	for _, r := range images {
		for _, v := range r.RepoTags {
			if v == image {
				hasImageExist = true
			}
		}
	}

	return hasImageExist, nil
}

func NewVirtualMachine(option *Options) (*VirtualMachine, error) {
	ctx := context.Background()

	cli, err := client.NewEnvClient()

	if err != nil {
		return nil, err
	}

	if exist, err := isImageExist(cli, ctx, option.Image); err != nil {
		return nil, err
	} else if !exist {
		cmd := exec.Command("docker", "pull", option.Image)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return nil, err
		}
	}

	containerWorkingDir := "/root/app"

	cfg := container.Config{
		Image:        option.Image,
		Tty:          true,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		OpenStdin:    true,
		WorkingDir:   containerWorkingDir,
	}

	if option.Commands != nil {
		cfg.Cmd = *option.Commands
	}

	resp, err := cli.ContainerCreate(ctx, &cfg, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: option.Cwd,
				Target: containerWorkingDir,
			},
		},
	}, nil, "")

	if err != nil {
		return nil, err
	}

	return &VirtualMachine{
		id:     resp.ID,
		ctx:    ctx,
		client: cli,
	}, nil
}

func (v *VirtualMachine) Start() error {
	if err := v.client.ContainerStart(v.ctx, v.id, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

// wait machine ready
func (v *VirtualMachine) Wait() error {
	_, err := v.client.ContainerWait(v.ctx, v.id)

	if err != nil {
		return err
	}

	return nil
}

// wait machine ready
func (v *VirtualMachine) Log() error {
	stdout, err := v.client.ContainerLogs(v.ctx, v.id, types.ContainerLogsOptions{ShowStdout: true})

	if err != nil {
		return err
	}

	stderr, err := v.client.ContainerLogs(v.ctx, v.id, types.ContainerLogsOptions{ShowStderr: true})

	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	wg.Add(1)

	go func() {
		_, _ = io.Copy(os.Stdout, stdout)
		wg.Done()
	}()

	go func() {
		_, _ = io.Copy(os.Stderr, stderr)
		wg.Done()
	}()

	wg.Wait()

	return nil
}

// wait machine ready
func (v *VirtualMachine) Attach() error {
	bashPath := "/bin/sh"

	if runtime.GOOS == "windows" {
		bashPath = "cmd.exe"
	}

	cmd := exec.Command("docker", "exec", "-it", v.id, bashPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// wait machine ready
func (v *VirtualMachine) Destroy() error {
	err := v.client.ContainerRemove(v.ctx, v.id, types.ContainerRemoveOptions{
		Force: true,
	})

	if err != nil {
		return err
	}

	return nil
}
