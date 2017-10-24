package evaluator

import (
	"fmt"

	"github.com/jolisper/monkey/ast"
	"github.com/jolisper/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch typedNode := node.(type) {
	case *ast.Program:
		return evalProgram(typedNode, env)

	case *ast.ExpressionStatement:
		return Eval(typedNode.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: typedNode.Value}

	case *ast.Boolean:
		return nativeBooleanToBooleanObject(typedNode.Value)

	case *ast.PrefixExpression:
		right := Eval(typedNode.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(typedNode.Operator, right)

	case *ast.InfixExpression:
		left := Eval(typedNode.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(typedNode.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(typedNode.Operator, left, right, env)

	case *ast.BlockStatement:
		return evalBlockStatement(typedNode, env)

	case *ast.IfExpression:
		return evalIfExpression(typedNode, env)

	case *ast.ReturnStatement:
		val := Eval(typedNode.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(typedNode.Value, env)
		if isError(val) {
			return val
		}

		env.Set(typedNode.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(typedNode, env)

	case *ast.FunctionLiteral:
		params := typedNode.Parameters
		body := typedNode.Body
		return &object.Function{
			Parameters: params,
			Env:        env,
			Body:       body,
		}
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBooleanToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object, env *object.Environment) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBooleanToBooleanObject(left == right)
	case operator == "!=":
		return nativeBooleanToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, rigth object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rigthValue := rigth.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rigthValue}
	case "-":
		return &object.Integer{Value: leftValue - rigthValue}
	case "*":
		return &object.Integer{Value: leftValue * rigthValue}
	case "/":
		return &object.Integer{Value: leftValue / rigthValue}
	case "<":
		return nativeBooleanToBooleanObject(leftValue < rigthValue)
	case ">":
		return nativeBooleanToBooleanObject(leftValue > rigthValue)
	case "==":
		return nativeBooleanToBooleanObject(leftValue == rigthValue)
	case "!=":
		return nativeBooleanToBooleanObject(leftValue != rigthValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, rigth.Type())
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
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}

	return val
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
