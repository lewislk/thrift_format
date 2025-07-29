package formater

import (
	"errors"
	"fmt"

	"github.com/cloudwego/thriftgo/parser"
	"github.com/lewislk/go-utils/files"
	"github.com/lewislk/thrift_format/executor"
	"github.com/lewislk/thrift_format/logs"
	"github.com/urfave/cli/v2"
)

type Formater struct {
	filePath        string
	lineSelectStart int64
	lineSelectEnd   int64
	overwrite       bool
}

func NewFormater(ctx *cli.Context) *Formater {
	logs.Verbose = ctx.Bool("verbose")
	return &Formater{
		filePath:        ctx.String("file"),
		lineSelectStart: ctx.Int64("line_start"),
		lineSelectEnd:   ctx.Int64("line_end"),
		overwrite:       ctx.Bool("overwrite"),
	}
}

func (f *Formater) DoFormat() error {
	thrift, err := parser.ParseFile(f.filePath, nil, false)
	if err != nil {
		return errors.New(fmt.Sprintf("parse thrift file err, '%s' is not valid thrift file", f.filePath))
	}
	lines, err := files.ReadFile(f.filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("read thrift file err:%v", err))
	}
	result, err := executor.NewExecutor(lines, thrift).Exec(f.lineSelectStart, f.lineSelectEnd)
	if err != nil {
		return err
	}
	if f.overwrite {
		err := files.OverwriteFile(f.filePath, result)
		if err != nil {
			return errors.New(fmt.Sprintf("write output err:%v", err))
		}
	} else {
		fmt.Println(result)
	}
	logs.Info("format success")
	return nil
}
