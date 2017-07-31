package evaluator

import (
	"github.com/jolisper/monkey/ast"
	"github.com/jolisper/monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	}

	return nil
}

func evalStatements(stms []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stms {
		result = Eval(statement)
	}

	return result
}
