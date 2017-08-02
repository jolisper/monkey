package evaluator

import (
	"github.com/jolisper/monkey/ast"
	"github.com/jolisper/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch typedNode := node.(type) {
	case *ast.Program:
		return evalStatements(typedNode.Statements)

	case *ast.ExpressionStatement:
		return Eval(typedNode.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: typedNode.Value}

	case *ast.Boolean:
		return nativeBooleanToBooleanObject(typedNode.Value)
	}

	return nil
}

func nativeBooleanToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func evalStatements(stms []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stms {
		result = Eval(statement)
	}

	return result
}
