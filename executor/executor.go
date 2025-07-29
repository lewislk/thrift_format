package executor

import (
	"strings"

	"github.com/cloudwego/thriftgo/parser"
	"github.com/lewislk/thrift_format/common"
	"github.com/lewislk/thrift_format/execution"
	"github.com/lewislk/thrift_format/factory"
	"github.com/lewislk/thrift_format/mctx"
)

type Executor struct {
	ctx *mctx.Context

	thrift *parser.Thrift

	curExecution execution.Execution
}

func NewExecutor(lines []string, thrift *parser.Thrift) *Executor {
	structMap := make(map[string]*parser.StructLike)
	enumMap := make(map[string]*parser.Enum)
	for _, s := range thrift.GetStructs() {
		structMap[s.Name] = s
	}
	for _, e := range thrift.GetEnums() {
		enumMap[e.Name] = e
	}
	return &Executor{
		ctx: &mctx.Context{
			Lines:     lines,
			StructMap: structMap,
			EnumMap:   enumMap,
			Constants: thrift.Constants,
		},
		thrift: thrift,
	}
}

func (e *Executor) Exec(startLine, endLine int64) (string, error) {
	result := make([]string, 0)
	for offset, line := range e.ctx.Lines {
		e.ctx.CurIdx = offset
		lineNum := int64(offset + 1)
		if startLine > 0 && endLine > 0 && endLine > startLine && (lineNum < startLine || lineNum > endLine) {
			result = append(result, line)
			continue
		}
		prefixType := common.GetPrefixType(line)
		if e.curExecution == nil {
			e.curExecution = factory.GetExecution(e.ctx, prefixType)
		} else {
			if !e.curExecution.CanContinue(prefixType) {
				e.curExecution = factory.GetExecution(e.ctx, prefixType)
			}
		}

		output := e.curExecution.Process(prefixType)
		result = append(result, strings.TrimRight(output, " "))
		if e.curExecution.IsFinish() {
			e.curExecution = nil
		}
	}
	return strings.Join(result, "\n"), nil
}
