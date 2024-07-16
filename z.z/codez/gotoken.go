package codez

import (
	"go/token"
)

// this code is taken from go/parser/parser.go
var stmtStart = map[token.Token]bool{
	token.BREAK:       true,
	token.CONST:       true,
	token.CONTINUE:    true,
	token.DEFER:       true,
	token.FALLTHROUGH: true,
	token.FOR:         true,
	token.GO:          true,
	token.GOTO:        true,
	token.IF:          true,
	token.RETURN:      true,
	token.SELECT:      true,
	token.SWITCH:      true,
	token.TYPE:        true,
	token.VAR:         true,
}

// this code is taken from go/parser/parser.go
var declStart = map[token.Token]bool{
	token.IMPORT: true,
	token.CONST:  true,
	token.TYPE:   true,
	token.VAR:    true,
}
