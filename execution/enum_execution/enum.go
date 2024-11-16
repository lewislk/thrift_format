package enum_execution

import (
	"fmt"
	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/parser/token"
	"github.com/lewis-buji/thrift_format/common"
	"github.com/lewis-buji/thrift_format/consts"
	"github.com/lewis-buji/thrift_format/execution"
	"github.com/lewis-buji/thrift_format/execution/base_execution"
	"github.com/lewis-buji/thrift_format/logs"
	"github.com/lewis-buji/thrift_format/mctx"
	"github.com/lewis-buji/thrift_format/utils"
)

type EnumExecution struct {
	*base_execution.BaseExecution

	format   string
	valueMap map[string]*parser.EnumValue
	annoMap  map[string]string
}

func NewEnumExecution(ctx *mctx.Context) execution.Execution {
	nameMaxLen, valueMaxLen, annoMaxLen := 0, 0, 0
	valueMap := make(map[string]*parser.EnumValue)
	annoMap := make(map[string]string)

	name := common.FindFirst(ctx.Lines[ctx.CurIdx], token.Identifier)
	curEnum := ctx.EnumMap[name]
	for _, value := range curEnum.Values {
		valueMap[value.Name] = value
		nameMaxLen = utils.Max(nameMaxLen, len(value.Name))
		valueMaxLen = utils.Max(valueMaxLen, len(fmt.Sprintf(`%d`, value.Value)))
		anno := common.GetAnnotation(value.Annotations)
		annoMaxLen = utils.Max(annoMaxLen, len(anno))
		annoMap[value.Name] = anno
	}

	// example
	// enum ExampleType {
	//     Type1 = 1        (tag.key='value') // comment
	//     NAME    VALUE       TAG               COMMENT
	//      %s       %d         %s                 %s
	// }
	format := fmt.Sprintf("    %%-%ds = %%-%dd %%-%ds %%s", nameMaxLen, valueMaxLen, annoMaxLen)

	ctx.Status = consts.InEnum
	return &EnumExecution{
		BaseExecution: &base_execution.BaseExecution{
			Ctx: ctx,
		},
		format:   format,
		valueMap: valueMap,
		annoMap:  annoMap,
	}
}

func (e *EnumExecution) CanContinue(prefixType token.Tok) bool {
	return !e.IsBlockType(prefixType) && e.Ctx.Status == consts.InEnum
}

func (e *EnumExecution) Process(prefixType token.Tok) string {
	line := e.Ctx.Lines[e.Ctx.CurIdx]
	output := line
	switch prefixType {
	case token.Identifier: // in field
		output = e.FormatLine(line)
	case token.RBrace: // exit struct
		e.Ctx.Status = consts.InOut
	default:
	}
	return output
}

func (e *EnumExecution) IsFinish() bool {
	return e.Ctx.Status != consts.InEnum
}

func (e *EnumExecution) FormatLine(line string) string {
	comment := common.FormatComment(line)
	valueName := common.FindFirst(line, token.Identifier)
	value, ok := e.valueMap[valueName]
	if !ok {
		logs.ErrorF(`line: '%s' can not find value`, line)
		return line
	}
	return fmt.Sprintf(e.format, value.Name, value.Value, e.annoMap[value.Name], comment)
}
