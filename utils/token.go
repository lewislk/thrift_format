package utils

import (
	"bytes"
	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/parser/token"
	"strings"
)

func GetPrefixType(line string) token.Tok {
	line = strings.TrimSpace(line)
	return token.NewTokenizer(bytes.NewBufferString(line)).Next().Tok
}

func FindFirst(line string, want token.Tok) string {
	tokenizer := token.NewTokenizer(bytes.NewBufferString(line))
	for {
		t := tokenizer.Next()
		if t.Tok == want {
			return t.AsString()
		}
		if t.Tok == token.EOF {
			return ""
		}
	}
}

func FindConst(line string, constants []*parser.Constant) *parser.Constant {
	for _, constant := range constants {
		for _, identifier := range findAllToken(line, token.Identifier) {
			if identifier == constant.GetName() {
				return constant
			}
		}
	}
	return nil
}

func findAllToken(line string, want token.Tok) (arr []string) {
	tokenizer := token.NewTokenizer(bytes.NewBufferString(line))
	for {
		t := tokenizer.Next()
		if t.Tok == want {
			arr = append(arr, t.AsString())
		}
		if t.Tok == token.EOF {
			return
		}
	}
}
