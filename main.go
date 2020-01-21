package main

import (
	"fmt"
	"os"

	"github.com/axetroy/dvs/internal/command"
	"github.com/axetroy/dvs/internal/dir"
	"github.com/axetroy/dvs/internal/version"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "dvs"
	app.Usage = "Docker-based Virtual System"
	app.Version = version.GetCurrentUsingVersion()
	app.Authors = []*cli.Author{
		{
			Name:  "Axetroy",
			Email: "axetroy.dev@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "image",
			Aliases: []string{"i"},
			Usage:   "Specifying the running image. eg. ubuntu:19.10",
			Value:   "alpine:3.11", // default image
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "run",
			Usage: "Run command in container",
			Action: func(c *cli.Context) error {
				return command.Run(c.Args().Slice(), &command.RunOption{
					Image: c.String("image"),
				})
			},
		},
		{
			Name:  "upgrade",
			Usage: "Upgrade dvs",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "force", Aliases: []string{"f"}, Usage: "Force upgrade", Value: false},
			},
			Action: func(c *cli.Context) error {
				return command.Upgrade(c.Args().First(), c.Bool("force"))
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		return command.Repl(&command.ReplOption{
			Image: c.String("image"),
		})
	}

	// regardless of the result, the cache directory should be delete
	if err := app.Run(os.Args); err != nil {
		if os.Getenv("DEBUG") != "" {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Println(err.Error())
			fmt.Printf("run with environment variables %s to print more information\n", color.GreenString("DEBUG=1"))
		}
		_ = os.RemoveAll(dir.CacheDir)
		os.Exit(1)
	} else {
		_ = os.RemoveAll(dir.CacheDir)
	}
}
