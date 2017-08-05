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

	case *ast.PrefixExpression:
		right := Eval(typedNode.Right)
		return evalPrefixExpression(typedNode.Operator, right)
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

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
