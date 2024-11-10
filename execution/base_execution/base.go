package base_execution

import (
	"gitee.com/liukunc9/thrift_format/mctx"
	"github.com/cloudwego/thriftgo/parser/token"
	"slices"
)

var blockTypeList = []token.Tok{
	token.Const,
	token.Enum,
	token.Struct,
}

type BaseExecution struct {
	Ctx *mctx.Context
}

func (e *BaseExecution) IsBlockType(prefixType token.Tok) bool {
	return slices.Contains(blockTypeList, prefixType)
}

func (e *BaseExecution) IsMatch(prefixType token.Tok) bool {
	panic("implement me")
}

func (e *BaseExecution) Process(prefixType token.Tok) string {
	panic("implement me")
}

func (e *BaseExecution) IsFinish() bool {
	panic("implement me")
}
