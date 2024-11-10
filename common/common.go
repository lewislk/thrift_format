package common

import (
	"bytes"
	"fmt"
	"gitee.com/liukunc9/go-utils/strs"
	"gitee.com/liukunc9/thrift_format/logs"
	"gitee.com/liukunc9/thrift_format/utils"
	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/parser/token"
	"strconv"
	"strings"
)

func ConvertConstValue2Str(t *parser.ConstValue) string {
	if t == nil {
		return ""
	}
	var val string
	switch t.Type {
	case parser.ConstType_ConstDouble:
		val = fmt.Sprintf("%f", *t.TypedValue.Double)
	case parser.ConstType_ConstInt:
		val = fmt.Sprintf("%d", *t.TypedValue.Int)
	case parser.ConstType_ConstLiteral:
		val = fmt.Sprintf("\"%s\"", *t.TypedValue.Literal)
	case parser.ConstType_ConstIdentifier:
		val = *t.TypedValue.Identifier
	case parser.ConstType_ConstList:
		if len(t.TypedValue.List) == 0 {
			return "{}"
		}
		var ss []string
		for _, item := range t.TypedValue.List {
			ss = append(ss, ConvertConstValue2Str(item))
		}
		val = fmt.Sprintf("[%s]", strings.Join(ss, ","))
	case parser.ConstType_ConstMap:
		if len(t.TypedValue.Map) == 0 {
			return "{}"
		}
		var ss []string
		for _, item := range t.TypedValue.Map {
			key := ConvertConstValue2Str(item.Key)
			val := ConvertConstValue2Str(item.Value)
			pair := key + ":" + val
			ss = append(ss, pair)
		}
		val = fmt.Sprintf("{%s}", strings.Join(ss, ","))
	default:
		return fmt.Sprintf("%+v", *t)
	}
	return val
}

func GetAnnotation(annotations []*parser.Annotation) string {
	if len(annotations) == 0 {
		return ""
	}
	buf := bytes.Buffer{}
	for _, annotation := range annotations {
		if len(annotation.GetValues()) == 0 {
			buf.WriteString(fmt.Sprintf("%s = ''", annotation.GetKey()))
			continue
		}
		val := annotation.GetValues()[0]
		if strings.Contains(val, "\"") {
			val = strings.ReplaceAll(val, "\"", "\\\"")
		}
		unquote, err := strconv.Unquote(fmt.Sprintf(`"%s"`, val))
		if err != nil {
			logs.Warn("Unquote error: %v", val)
			buf.WriteString(fmt.Sprintf(`%s = '%s',`, annotation.GetKey(), val))
			continue
		}
		buf.WriteString(fmt.Sprintf(`%s = '%s',`, annotation.GetKey(), strconv.Quote(unquote)))
	}
	buf.Truncate(buf.Len() - 1)
	return fmt.Sprintf(`(%s)`, buf.String())
}

func FormatComment(line string) string {
	var comment string
	if comment = utils.FindFirst(line, token.LineComment); !strs.IsEmpty(comment) {
		idx := strings.Index(comment, "//")
		return fmt.Sprintf(`// %s`, strings.TrimSpace(comment[idx+2:]))
	} else if comment = utils.FindFirst(line, token.UnixComment); !strs.IsEmpty(comment) {
		idx := strings.Index(comment, "#")
		return fmt.Sprintf(`# %s`, strings.TrimSpace(comment[idx+1:]))
	}

	return ""
}
