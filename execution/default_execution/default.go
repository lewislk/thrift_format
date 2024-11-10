package default_execution

import (
	"gitee.com/liukunc9/thrift_format/consts"
	"gitee.com/liukunc9/thrift_format/execution"
	"gitee.com/liukunc9/thrift_format/execution/base_execution"
	"gitee.com/liukunc9/thrift_format/mctx"
	"github.com/cloudwego/thriftgo/parser/token"
)

type DefaultExecution struct {
	*base_execution.BaseExecution
}

func NewDefaultExecution(ctx *mctx.Context) execution.Execution {
	ctx.Status = consts.InOut
	return &DefaultExecution{
		BaseExecution: &base_execution.BaseExecution{
			Ctx: ctx,
		},
	}
}

func (e *DefaultExecution) IsMatch(prefixType token.Tok) bool {
	return e.Ctx.Status == consts.InOut
}

func (e *DefaultExecution) Process(prefixType token.Tok) string {
	return e.Ctx.Lines[e.Ctx.CurIdx]
}

func (e *DefaultExecution) IsFinish() bool {
	return e.Ctx.Status != consts.InOut
}
