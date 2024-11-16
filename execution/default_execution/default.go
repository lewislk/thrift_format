package default_execution

import (
	"github.com/cloudwego/thriftgo/parser/token"
	"github.com/liukunc9/thrift_format/consts"
	"github.com/liukunc9/thrift_format/execution"
	"github.com/liukunc9/thrift_format/execution/base_execution"
	"github.com/liukunc9/thrift_format/mctx"
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

func (e *DefaultExecution) CanContinue(prefixType token.Tok) bool {
	return !e.IsBlockType(prefixType) && e.Ctx.Status == consts.InOut
}

func (e *DefaultExecution) Process(prefixType token.Tok) string {
	return e.Ctx.Lines[e.Ctx.CurIdx]
}

func (e *DefaultExecution) IsFinish() bool {
	return e.Ctx.Status != consts.InOut
}
