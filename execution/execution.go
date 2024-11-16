package execution

import "github.com/cloudwego/thriftgo/parser/token"

type Execution interface {
	CanContinue(prefixType token.Tok) bool
	Process(prefixType token.Tok) string
	IsFinish() bool
}
