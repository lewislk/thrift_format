package factory

import (
	"github.com/cloudwego/thriftgo/parser/token"
	"github.com/lewislk/thrift_format/execution"
	"github.com/lewislk/thrift_format/execution/const_execution"
	"github.com/lewislk/thrift_format/execution/default_execution"
	"github.com/lewislk/thrift_format/execution/enum_execution"
	"github.com/lewislk/thrift_format/execution/struct_execution"
	"github.com/lewislk/thrift_format/mctx"
)

type Constructor = func(ctx *mctx.Context) execution.Execution

var executionMap = map[token.Tok]Constructor{
	token.Struct: struct_execution.NewStructExecution,
	token.Enum:   enum_execution.NewEnumExecution,
	token.Const:  const_execution.NewConstExecution,
}

func GetExecution(ctx *mctx.Context, prefixType token.Tok) execution.Execution {
	cons := executionMap[prefixType]
	if cons == nil {
		cons = default_execution.NewDefaultExecution
	}
	return cons(ctx)
}
