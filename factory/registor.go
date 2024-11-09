package factory

import (
	"gitee.com/liukunc9/thrift_format/execution"
	"gitee.com/liukunc9/thrift_format/execution/const_execution"
	"gitee.com/liukunc9/thrift_format/execution/enum_execution"
	"gitee.com/liukunc9/thrift_format/execution/struct_execution"
	"gitee.com/liukunc9/thrift_format/mctx"
	"github.com/cloudwego/thriftgo/parser/token"
)

type Constructor = func(ctx *mctx.Context) execution.Execution

var ExecutionMap = map[token.Tok]Constructor{
	token.Struct: struct_execution.NewStructExecution,
	token.Enum:   enum_execution.NewEnumExecution,
	token.Const:  const_execution.NewConstExecution,
}
