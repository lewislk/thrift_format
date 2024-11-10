package execution

import "github.com/cloudwego/thriftgo/parser/token"

type Execution interface {
	IsMatch(prefixType token.Tok) bool
	Process(prefixType token.Tok) string
	IsFinish() bool
}
