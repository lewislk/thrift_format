package enum_execution

import (
	"gitee.com/liukunc9/thrift_format/consts"
	"gitee.com/liukunc9/thrift_format/execution"
	"gitee.com/liukunc9/thrift_format/execution/base_execution"
	"gitee.com/liukunc9/thrift_format/mctx"
)

type EnumExecution struct {
	*base_execution.BaseExecution
}

func NewEnumExecution(ctx *mctx.Context) execution.Execution {
	ctx.Status = consts.InEnum
	return &EnumExecution{
		BaseExecution: &base_execution.BaseExecution{
			Ctx: ctx,
		},
	}
}
