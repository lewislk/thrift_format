package mctx

import (
	"github.com/cloudwego/thriftgo/parser"
	"github.com/liukunc9/thrift_format/consts"
)

type Context struct {
	Lines     []string
	StructMap map[string]*parser.StructLike
	EnumMap   map[string]*parser.Enum
	Constants []*parser.Constant

	Status consts.Status
	CurIdx int
}
