package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type RunOption struct {
	Image string // 运行的镜像
}

type Image struct {
	Name string
	Tag  string
}

func ImageParser(image string) Image {
	arr := strings.Split(image, ":")

	name := arr[0]
	var tag string

	if len(arr) > 1 {
		tag = arr[1]
	} else {
		tag = "latest"
	}

	return Image{
		Name: name,
		Tag:  tag,
	}
}

func Run(command []string, option *RunOption) error {
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
		_, err := cli.ImagePull(ctx, "docker.io/library/"+option.Image, types.ImagePullOptions{})

		if err != nil {
			return err
		}
	}

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	containerWorkingDir := "/root/app"

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      option.Image,
		Cmd:        command,
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
		_ = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			Force: true,
		})
	}()

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	// wait container ready
	_, err = cli.ContainerWait(ctx, resp.ID)

	if err != nil {
		return err
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})

	if err != nil {
		panic(err)
	}

	// TODO: use
	// stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	_, err = io.Copy(os.Stdout, out)

	if err != nil {
		return err
	}

	return nil
}
