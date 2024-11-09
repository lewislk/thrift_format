package struct_execution

import (
	"gitee.com/liukunc9/thrift_format/consts"
	"gitee.com/liukunc9/thrift_format/execution"
	"gitee.com/liukunc9/thrift_format/execution/base_execution"
	"gitee.com/liukunc9/thrift_format/mctx"
)

type StructExecution struct {
	*base_execution.BaseExecution
}

func NewStructExecution(ctx *mctx.Context) execution.Execution {
	ctx.Status = consts.InStruct
	return &StructExecution{
		BaseExecution: &base_execution.BaseExecution{
			Ctx: ctx,
		},
	}
}
