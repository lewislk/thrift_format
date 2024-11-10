package base_execution

import (
	"gitee.com/liukunc9/thrift_format/mctx"
	"github.com/cloudwego/thriftgo/parser/token"
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
	for _, blockType := range blockTypeList {
		if prefixType == blockType {
			return true
		}
	}
	return false
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
