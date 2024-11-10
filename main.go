package main

import (
	"gitee.com/liukunc9/thrift_format/formater"
	"gitee.com/liukunc9/thrift_format/logs"
	"github.com/urfave/cli/v2"
	"os"
	"runtime/debug"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logs.ErrorF("panic: %v", r)
			logs.DebugF("%s", string(debug.Stack()))
		}
	}()
	app := &cli.App{
		Name:    "thrift_format",
		Usage:   "thrift_format -f `FilePath`",
		Action:  action,
		Version: "0.0.3",
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
			&cli.Int64Flag{
				Name:     "line_start",
				Aliases:  []string{"ls"},
				Usage:    "line select start",
				Required: false,
			},
			&cli.Int64Flag{
				Name:     "line_end",
				Aliases:  []string{"le"},
				Usage:    "line select end",
				Required: false,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				DefaultText: "false",
				Usage:       "print debug log",
				Required:    false,
				Value:       false,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		logs.ErrorF("%v", err)
	}
}

func action(ctx *cli.Context) error {
	f := formater.NewFormater(ctx)
	return f.DoFormat()
}
