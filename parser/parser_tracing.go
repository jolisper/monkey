package parser

import (
	"fmt"
	"os"
	"strings"
)

var traceLevel int = 0

const traceIndentPlaceholder string = "\t"

func indentLevel() string {
	return strings.Repeat(traceIndentPlaceholder, traceLevel-1)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", indentLevel(), fs)
}

func incIdent() { traceLevel = traceLevel + 1 }
func decIdent() { traceLevel = traceLevel - 1 }

func trace(msg string) string {
	if os.Getenv("MONKEY_PARSER_TRACING_ON") != "" {
		incIdent()
		tracePrint("BEGIN " + msg)
	}
	return msg
}

func untrace(msg string) {
	if os.Getenv("MONKEY_PARSER_TRACING_ON") != "" {
		tracePrint("END " + msg)
		decIdent()
	}
}
