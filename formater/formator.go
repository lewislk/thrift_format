package formater

import (
	"fmt"
	"gitee.com/liukunc9/go-utils/files"
	"gitee.com/liukunc9/thrift_format/executor"
	"gitee.com/liukunc9/thrift_format/logs"
	"github.com/cloudwego/thriftgo/parser"
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
		logs.Error("parse thrift file err:%v", err)
		return err
	}
	lines, err := files.ReadFile(f.filePath)
	if err != nil {
		logs.Error("read thrift file err:%v", err)
	}
	result, err := executor.NewExecutor(lines, thrift).Exec(f.lineSelectStart, f.lineSelectEnd)
	if err != nil {
		return err
	}
	if f.overwrite {
		err := files.OverwriteFile(f.filePath, result)
		if err != nil {
			logs.Error("write output err:%v", err)
			return err
		}
	} else {
		fmt.Println(result)
	}
	logs.Info("format success")
	return nil
}
