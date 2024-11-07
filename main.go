package main

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"os"
	"runtime/debug"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			color.Red("panic: %v", r)
			color.Red("%s", string(debug.Stack()))
		}
	}()
	app := &cli.App{
		Name:  "thrift_format",
		Usage: "format thrift file",
		Action: func(c *cli.Context) error {
			return nil
		},
		Version: "0.1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "`file` to format",
				Required: false,
			},
			&cli.BoolFlag{
				Name:        "overwrite",
				Aliases:     []string{"o"},
				DefaultText: "true",
				Usage:       "overwrite file or not, default true",
				Required:    false,
				Value:       true,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				DefaultText: "false",
				Usage:       "display debug log",
				Required:    false,
				Value:       true,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		color.Red("%v", err)
	}
}
