package const_execution

import (
	"gitee.com/liukunc9/thrift_format/consts"
	"gitee.com/liukunc9/thrift_format/execution"
	"gitee.com/liukunc9/thrift_format/execution/base_execution"
	"gitee.com/liukunc9/thrift_format/mctx"
)

type ConstExecution struct {
	*base_execution.BaseExecution
}

func NewConstExecution(ctx *mctx.Context) execution.Execution {
	ctx.Status = consts.InConst
	return &ConstExecution{
		BaseExecution: &base_execution.BaseExecution{
			Ctx: ctx,
		},
	}
}
