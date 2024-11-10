package const_execution

import (
	"fmt"
	"gitee.com/liukunc9/thrift_format/common"
	"gitee.com/liukunc9/thrift_format/consts"
	"gitee.com/liukunc9/thrift_format/execution"
	"gitee.com/liukunc9/thrift_format/execution/base_execution"
	"gitee.com/liukunc9/thrift_format/logs"
	"gitee.com/liukunc9/thrift_format/mctx"
	"gitee.com/liukunc9/thrift_format/utils"
	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/parser/token"
	"strings"
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
		if utils.GetPrefixType(currentLine) != token.Const {
			break
		}
		currentConst := utils.FindConst(currentLine, ctx.Constants)
		if currentConst == nil {
			continue
		}
		constMap[currentLine] = currentConst
		typeStr := fmt.Sprintf("%v", currentConst.Type)
		value := common.ConvertConstValue2Str(currentConst.Value)
		valueMap[currentConst.Name] = value
		typeMaxLen = max(typeMaxLen, len(typeStr))
		nameMaxLen = max(nameMaxLen, len(currentConst.Name))
		valueMaxLen = max(valueMaxLen, getValueLen(value))
		anno := common.GetAnnotation(currentConst.GetAnnotations())
		annoMaxLen = max(annoMaxLen, len(anno))
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

func (e *ConstExecution) IsMatch(prefixType token.Tok) bool {
	return !e.IsBlockType(prefixType) && e.Ctx.Status == consts.InConst
}

func (e *ConstExecution) Process(prefixType token.Tok) string {
	line := e.Ctx.Lines[e.Ctx.CurIdx]
	output := line
	switch prefixType {
	case token.IntLiteral: // in field
		output = e.FormatLine(line)
	case token.RBrace: // exit struct
		e.Ctx.Status = consts.InOut
	default:
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
		logs.Error(`line: '%v' can not find cost value`, line)
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
