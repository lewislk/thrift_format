package mctx

import (
	"gitee.com/liukunc9/thrift_format/consts"
	"github.com/cloudwego/thriftgo/parser"
)

type Context struct {
	Status    consts.Status
	StructMap map[string]*parser.StructLike
	EnumMap   map[string]*parser.Enum
}
