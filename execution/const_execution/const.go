package const_execution

import (
	"fmt"
	"strings"

	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/parser/token"
	"github.com/lewislk/thrift_format/common"
	"github.com/lewislk/thrift_format/consts"
	"github.com/lewislk/thrift_format/execution"
	"github.com/lewislk/thrift_format/execution/base_execution"
	"github.com/lewislk/thrift_format/logs"
	"github.com/lewislk/thrift_format/mctx"
	"github.com/lewislk/thrift_format/utils"
)

type ConstExecution struct {
	*base_execution.BaseExecution

	format   string
	valueMap map[string]string
	constMap map[string]*parser.Constant
	annoMap  map[string]string
}

func NewConstExecution(ctx *mctx.Context) execution.Execution {
	typeMaxLen, nameMaxLen, valueMaxLen, annoMaxLen := 0, 0, 0, 0
	valueMap := make(map[string]string)
	constMap := make(map[string]*parser.Constant)
	annoMap := make(map[string]string)

	for idx := ctx.CurIdx; idx < len(ctx.Lines); idx++ {
		currentLine := ctx.Lines[idx]
		if common.GetPrefixType(currentLine) != token.Const {
			break
		}
		currentConst := common.FindConst(currentLine, ctx.Constants)
		if currentConst == nil {
			continue
		}
		constMap[currentLine] = currentConst
		typeStr := fmt.Sprintf("%v", currentConst.Type)
		value := common.ConvertConstValue2Str(currentConst.Value)
		valueMap[currentConst.Name] = value
		typeMaxLen = utils.Max(typeMaxLen, len(typeStr))
		nameMaxLen = utils.Max(nameMaxLen, len(currentConst.Name))
		valueMaxLen = utils.Max(valueMaxLen, getValueLen(value))
		anno := common.GetAnnotation(currentConst.GetAnnotations())
		annoMaxLen = utils.Max(annoMaxLen, len(anno))
		annoMap[currentConst.Name] = anno
	}
	format := fmt.Sprintf("const %%-%ds %%-%ds = %%-%ds %%-%ds %%s",
		typeMaxLen, nameMaxLen, valueMaxLen, annoMaxLen)

	ctx.Status = consts.InConst
	return &ConstExecution{
		BaseExecution: &base_execution.BaseExecution{
			Ctx: ctx,
		},
		format:   format,
		valueMap: valueMap,
		constMap: constMap,
		annoMap:  annoMap,
	}
}

func (e *ConstExecution) CanContinue(prefixType token.Tok) bool {
	return prefixType == token.Const && e.Ctx.Status == consts.InConst
}

func (e *ConstExecution) Process(prefixType token.Tok) string {
	line := e.Ctx.Lines[e.Ctx.CurIdx]
	output := line
	switch prefixType {
	case token.Const: // in field
		output = e.FormatLine(line)
	default:
		e.Ctx.Status = consts.InOut
	}
	return output
}

func (e *ConstExecution) IsFinish() bool {
	return e.Ctx.Status != consts.InConst
}

func (e *ConstExecution) FormatLine(line string) string {
	comment := common.FormatComment(line)
	c, ok := e.constMap[line]
	if !ok {
		logs.ErrorF(`line: '%v' can not find cost value`, line)
		return line
	}
	value := e.valueMap[c.Name]
	return fmt.Sprintf(e.format, c.Type, c.Name, value, e.annoMap[c.Name], comment)
}

func getValueLen(value string) int {
	idx := strings.Index(value, "\n")
	if idx == -1 {
		return len(value)
	}
	return idx + 1
}
