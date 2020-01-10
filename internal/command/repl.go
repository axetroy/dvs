package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type ReplOption struct {
	Image string
}

func Repl(option *ReplOption) error {
	ctx := context.Background()

	cli, err := client.NewEnvClient()

	if err != nil {
		return err
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{
		All: true,
	})

	if err != nil {
		return err
	}

	hasImageExist := false

	targetImg := ImageParser(option.Image)

	for _, r := range images {
		for _, v := range r.RepoTags {
			img := ImageParser(v)

			if img.Name == targetImg.Name && img.Tag == targetImg.Tag {
				hasImageExist = true
			}
		}
	}

	if hasImageExist == false {
		fmt.Printf("Pulling image %v\n", option.Image)
		_, err := cli.ImagePull(ctx, "docker.io/library/"+option.Image, types.ImagePullOptions{
			All: true,
		})

		if err != nil {
			return err
		}
	}

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	containerWorkingDir := "/root/app"

	sandbox, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      option.Image,
		Tty:        true,
		OpenStdin:  true,
		WorkingDir: containerWorkingDir,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: cwd,
				Target: containerWorkingDir,
			},
		},
	}, nil, "")

	if err != nil {
		return err
	}

	// remove container
	defer func() {
		_ = cli.ContainerRemove(ctx, sandbox.ID, types.ContainerRemoveOptions{
			Force: true,
		})
	}()

	if err := cli.ContainerStart(ctx, sandbox.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	// TODO: not compatible with Window
	cmd := exec.Command("docker", "exec", "-it", sandbox.ID, "/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
