package struct_execution

import (
	"fmt"
	"gitee.com/liukunc9/go-utils/conv"
	"gitee.com/liukunc9/thrift_format/common"
	"gitee.com/liukunc9/thrift_format/consts"
	"gitee.com/liukunc9/thrift_format/execution"
	"gitee.com/liukunc9/thrift_format/execution/base_execution"
	"gitee.com/liukunc9/thrift_format/logs"
	"gitee.com/liukunc9/thrift_format/mctx"
	"gitee.com/liukunc9/thrift_format/utils"
	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/parser/token"
)

type StructExecution struct {
	*base_execution.BaseExecution

	fieldMap   map[int32]*parser.Field
	annoMap    map[int32]string
	defaultMap map[int32]string

	idMaxLen, requireTypeMaxLen, typeMaxLen, nameMaxLen, annotationMaxLen, defaultMaxLen int
}

func NewStructExecution(ctx *mctx.Context) execution.Execution {
	fieldMap := make(map[int32]*parser.Field)
	annoMap := make(map[int32]string)
	defaultMap := make(map[int32]string)
	idMaxLen, requireTypeMaxLen, typeMaxLen, nameMaxLen, annotationMaxLen, defaultMaxLen := 0, 0, 0, 0, 0, 0
	name := utils.FindFirst(ctx.Lines[ctx.CurIdx], token.Identifier)
	curStruct := ctx.StructMap[name]
	for _, field := range curStruct.Fields {
		fieldMap[field.ID] = field
		idMaxLen = max(idMaxLen, len(fmt.Sprintf("%d", field.ID)))
		requireTypeMaxLen = max(requireTypeMaxLen, len(getRequireType(field.GetRequiredness())))
		typeMaxLen = max(typeMaxLen, len(fmt.Sprintf("%v", field.Type)))
		nameMaxLen = max(nameMaxLen, len(field.Name))
		// default value
		defaultValue := getDefaultValue(field)
		defaultMap[field.ID] = defaultValue
		defaultMaxLen = max(defaultMaxLen, len(defaultValue))
		// annotation
		annotation := common.GetAnnotation(field.Annotations)
		annoMap[field.GetID()] = annotation
		annotationMaxLen = max(annotationMaxLen, len(annotation))
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
	}
}

func (e *StructExecution) IsMatch(prefixType token.Tok) bool {
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
	fieldID := int32(conv.Str2Int64(utils.FindFirst(line, token.IntLiteral), 0))
	field, ok := e.fieldMap[fieldID]
	if !ok {
		logs.ErrorF(`line: '%s' can not find field of id: %v`, line, field)
		return line
	}
	comment := common.FormatComment(line)
	requireType := getRequireType(field.Requiredness)

	var output string
	if len(requireType) > 0 {
		output = fmt.Sprintf(e.dynamicFormat(field, requireType), field.ID, requireType, field.Type, field.Name,
			e.defaultMap[fieldID], e.annoMap[fieldID], comment)
	} else {
		output = fmt.Sprintf(e.dynamicFormat(field, requireType), field.ID, field.Type, field.Name,
			e.defaultMap[fieldID], e.annoMap[fieldID], comment)
	}
	return output
}

func (e *StructExecution) IsFinish() bool {
	return e.Ctx.Status != consts.InStruct
}

func (e *StructExecution) dynamicFormat(field *parser.Field, requireType string) string {
	// 计算参数名称之前的最大长度。示例: 8: optional string score_url,即计算"8: optional string"的长度
	maxLenBeforeParamName := e.idMaxLen + 1 + 1 + e.typeMaxLen
	if e.requireTypeMaxLen > 0 {
		maxLenBeforeParamName += e.requireTypeMaxLen + 1
	}
	// 将所有不足长度的，补充空格到类型后面
	requireLen := len(requireType)
	if len(requireType) > 0 {
		// 计算 "8: optional "的长度
		typeFormatLen := maxLenBeforeParamName - (len(fmt.Sprintf("%d", field.ID)) + 1 + 1 + requireLen + 1)
		return fmt.Sprintf("    %%d: %%s %%-%dv %%-%ds%%-%ds %%-%ds %%s",
			typeFormatLen, e.nameMaxLen, e.defaultMaxLen, e.annotationMaxLen)
	} else {
		// 计算"8: "的长度
		typeFormatLen := maxLenBeforeParamName - (len(fmt.Sprintf("%d", field.ID)) + 1 + 1)
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
