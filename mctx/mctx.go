package mctx

import (
	"gitee.com/liukunc9/thrift_format/consts"
	"github.com/cloudwego/thriftgo/parser"
)

type Context struct {
	Lines     []string
	StructMap map[string]*parser.StructLike
	EnumMap   map[string]*parser.Enum
	Constants []*parser.Constant

	Status consts.Status
	CurIdx int
}
