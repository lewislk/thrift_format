package struct_execution

import (
	"fmt"
	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/parser/token"
	"github.com/liukunc9/go-utils/conv"
	"github.com/liukunc9/thrift_format/common"
	"github.com/liukunc9/thrift_format/consts"
	"github.com/liukunc9/thrift_format/execution"
	"github.com/liukunc9/thrift_format/execution/base_execution"
	"github.com/liukunc9/thrift_format/logs"
	"github.com/liukunc9/thrift_format/mctx"
	"github.com/liukunc9/thrift_format/utils"
)

type StructExecution struct {
	*base_execution.BaseExecution

	fieldMap   map[int32]*parser.Field
	annoMap    map[int32]string
	defaultMap map[int32]string

	idMaxLen, requireTypeMaxLen, typeMaxLen, nameMaxLen, annotationMaxLen, defaultMaxLen, paramBeforeMaxLen int
}

func NewStructExecution(ctx *mctx.Context) execution.Execution {
	fieldMap := make(map[int32]*parser.Field)
	annoMap := make(map[int32]string)
	defaultMap := make(map[int32]string)
	idMaxLen, requireTypeMaxLen, typeMaxLen, nameMaxLen, annotationMaxLen, defaultMaxLen, paramBeforeMaxLen := 0, 0, 0, 0, 0, 0, 0
	name := common.FindFirst(ctx.Lines[ctx.CurIdx], token.Identifier)
	curStruct := ctx.StructMap[name]
	for _, field := range curStruct.Fields {
		fieldMap[field.ID] = field
		idLen := len(fmt.Sprintf("%d", field.ID))
		idMaxLen = utils.Max(idMaxLen, idLen)
		requireTypeLen := len(getRequireType(field.GetRequiredness()))
		requireTypeMaxLen = utils.Max(requireTypeMaxLen, requireTypeLen)
		typeLen := len(fmt.Sprintf("%v", field.Type))
		typeMaxLen = utils.Max(typeMaxLen, typeLen)
		nameMaxLen = utils.Max(nameMaxLen, len(field.Name))
		if requireTypeLen > 0 {
			paramBeforeLen := idLen + 2 + requireTypeLen + 1 + typeLen
			paramBeforeMaxLen = utils.Max(paramBeforeMaxLen, paramBeforeLen)
		} else {
			paramBeforeLen := idLen + 2 + typeLen
			paramBeforeMaxLen = utils.Max(paramBeforeMaxLen, paramBeforeLen)
		}
		// default value
		defaultValue := getDefaultValue(field)
		defaultMap[field.ID] = defaultValue
		defaultMaxLen = utils.Max(defaultMaxLen, len(defaultValue))
		// annotation
		annotation := common.GetAnnotation(field.Annotations)
		annoMap[field.GetID()] = annotation
		annotationMaxLen = utils.Max(annotationMaxLen, len(annotation))
	}

	ctx.Status = consts.InStruct
	return &StructExecution{
		BaseExecution: &base_execution.BaseExecution{
			Ctx: ctx,
		},
		fieldMap:          fieldMap,
		annoMap:           annoMap,
		defaultMap:        defaultMap,
		idMaxLen:          idMaxLen,
		requireTypeMaxLen: requireTypeMaxLen,
		typeMaxLen:        typeMaxLen,
		nameMaxLen:        nameMaxLen,
		annotationMaxLen:  annotationMaxLen,
		defaultMaxLen:     defaultMaxLen,
		paramBeforeMaxLen: paramBeforeMaxLen,
	}
}

func (e *StructExecution) CanContinue(prefixType token.Tok) bool {
	return !e.IsBlockType(prefixType) && e.Ctx.Status == consts.InStruct
}

func (e *StructExecution) Process(prefixType token.Tok) string {
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

func (e *StructExecution) FormatLine(line string) string {
	fieldID := int32(conv.Str2Int64(common.FindFirst(line, token.IntLiteral), 0))
	field, ok := e.fieldMap[fieldID]
	if !ok {
		logs.ErrorF(`line: '%s' can not find field of id: %v`, line, field)
		return line
	}
	comment := common.FormatComment(line)
	requireType := getRequireType(field.Requiredness)

	var output string
	format := e.dynamicFormat(field, requireType)
	if len(requireType) > 0 {
		output = fmt.Sprintf(format, field.ID, requireType, field.Type, field.Name,
			e.defaultMap[fieldID], e.annoMap[fieldID], comment)
	} else {
		output = fmt.Sprintf(format, field.ID, field.Type, field.Name,
			e.defaultMap[fieldID], e.annoMap[fieldID], comment)
	}
	return output
}

func (e *StructExecution) IsFinish() bool {
	return e.Ctx.Status != consts.InStruct
}

func (e *StructExecution) dynamicFormat(field *parser.Field, requireType string) string {
	// 将所有不足长度的，补充空格到类型后面
	requireLen := len(requireType)
	if len(requireType) > 0 {
		// 计算 "8: optional "的长度
		typeFormatLen := e.paramBeforeMaxLen - (len(fmt.Sprintf("%d", field.ID)) + 1 + 1 + requireLen + 1)
		return fmt.Sprintf("    %%d: %%s %%-%dv %%-%ds%%-%ds %%-%ds %%s",
			typeFormatLen, e.nameMaxLen, e.defaultMaxLen, e.annotationMaxLen)
	} else {
		// 计算"8: "的长度
		typeFormatLen := e.paramBeforeMaxLen - (len(fmt.Sprintf("%d", field.ID)) + 1 + 1)
		return fmt.Sprintf("    %%d: %%-%dv %%-%ds%%-%ds %%-%ds %%s",
			typeFormatLen, e.nameMaxLen, e.defaultMaxLen, e.annotationMaxLen)
	}
}

func getRequireType(t parser.FieldType) string {
	switch t {
	case parser.FieldType_Default:
		return ""
	case parser.FieldType_Optional:
		return "optional"
	case parser.FieldType_Required:
		return "required"
	default:
		return ""
	}
}

func getDefaultValue(field *parser.Field) string {
	if field.Default == nil {
		return ""
	}
	return fmt.Sprintf(" = %v", common.ConvertConstValue2Str(field.Default))
}
