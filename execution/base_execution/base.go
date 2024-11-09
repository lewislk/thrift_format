package base_execution

import (
	"gitee.com/liukunc9/thrift_format/mctx"
	"github.com/cloudwego/thriftgo/parser/token"
)

type BaseExecution struct {
	Ctx *mctx.Context
}

func (e *BaseExecution) IsMatch(prefixType token.Tok) bool {
	//TODO implement me
	panic("implement me")
}

func (e *BaseExecution) Process(line string) string {
	//TODO implement me
	panic("implement me")
}

func (e *BaseExecution) IsFinish() bool {
	//TODO implement me
	panic("implement me")
}
