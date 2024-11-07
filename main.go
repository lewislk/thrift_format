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
		Name:    "thrift_format",
		Usage:   "format thrift file",
		Action:  action,
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "`file` to format",
				Required: true,
			},
			&cli.BoolFlag{
				Name:        "overwrite",
				Aliases:     []string{"o"},
				DefaultText: "true",
				Usage:       "overwrite file or not",
				Required:    false,
				Value:       true,
			},
			&cli.StringFlag{
				Name:     "lines",
				Aliases:  []string{"l"},
				Usage:    "`lines` selected to format",
				Required: false,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		color.Red("%v", err)
	}
}

func action(ctx *cli.Context) error {
	filePath := ctx.String("file")
	color.Green("file path: %s", filePath)
	lines := ctx.String("lines")
	color.Green("lines: %s", lines)
	return nil
}
