package projection

import (
	"goselect/parser/context"
)

type Expressions struct {
	expressions []*Expression
}

type expressionType int

const (
	ExpressionTypeValue     = 0
	ExpressionTypeAttribute = 1
	ExpressionTypeFunction  = 2
)

type Expression struct {
	eType     expressionType
	value     string
	attribute string
	function  *Function
}

type Function struct {
	name string
	args []*Expression
}

func expressionWithAttribute(attribute string) *Expression {
	return &Expression{
		eType:     ExpressionTypeAttribute,
		attribute: attribute,
	}
}

func expressionsWithAttributes(attributes []string) []*Expression {
	var expressions = make([]*Expression, len(attributes))
	for index, attribute := range attributes {
		expressions[index] = expressionWithAttribute(attribute)
	}
	return expressions
}

func expressionWithFunction(fn *Function) *Expression {
	return &Expression{
		eType:    ExpressionTypeFunction,
		function: fn,
	}
}

func expressionWithValue(value string) *Expression {
	return &Expression{
		eType: ExpressionTypeValue,
		value: value,
	}
}

func (expressions Expressions) count() int {
	return len(expressions.expressions)
}

func (expressions Expressions) displayableAttributes() []string {
	var functionAsString func(expression *Expression) string
	functionAsString = func(expression *Expression) string {
		if !expression.isAFunction() {
			if expression.eType == ExpressionTypeAttribute {
				return expression.attribute
			}
			return expression.value
		}
		var result = expression.function.name + "("
		for _, arg := range expression.function.args {
			result = result + functionAsString(arg) + ","
		}
		if len(expression.function.args) > 0 {
			result = result[0:len(result)-1] + ")"
		} else {
			result = result + ")"
		}
		return result
	}

	var attributes []string
	for _, expression := range expressions.expressions {
		if len(expression.attribute) != 0 {
			attributes = append(attributes, expression.attribute)
		} else if len(expression.value) != 0 {
			attributes = append(attributes, expression.value)
		} else {
			attributes = append(attributes, functionAsString(expression))
		}
	}
	return attributes
}

func (expressions Expressions) evaluateWith(fileAttributes *context.FileAttributes, functions *context.AllFunctions) ([]context.Value, error) {
	var execute func(expression *Expression) (context.Value, error)

	execute = func(expression *Expression) (context.Value, error) {
		if !expression.isAFunction() {
			return expression.getNonFunctionValue(fileAttributes), nil
		}
		var values []context.Value
		for _, arg := range expression.function.args {
			v, err := execute(arg)
			if err != nil {
				return context.EmptyValue(), nil
			}
			values = append(values, v)
		}
		return functions.Execute(expression.function.name, values...)
	}

	var values []context.Value
	for _, expression := range expressions.expressions {
		if !expression.isAFunction() {
			values = append(values, expression.getNonFunctionValue(fileAttributes))
		} else {
			value, err := execute(expression)
			if err != nil {
				return nil, err
			}
			values = append(values, value)
		}
	}
	return values, nil
}

func (expression Expression) isAFunction() bool {
	return expression.function != nil
}

func (expression Expression) getNonFunctionValue(fileAttributes *context.FileAttributes) context.Value {
	if expression.eType == ExpressionTypeAttribute {
		return fileAttributes.Get(expression.attribute)
	}
	return context.StringValue(expression.value)
}
