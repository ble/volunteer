package parse

import (
	"errors"
	"strconv"
)

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\r'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

var operationForChar map[byte]Operation

func trimWhitespace(in []byte) (*expr, []byte, error) {
	for len(in) > 0 && isSpace(in[0]) {
		in = in[1:]
	}
	return nil, in, nil
}

func Parse(input []byte) (Expr, error) {
	result, left, err := parseExpr(input)

	if err != nil {
		return nil, err
	}
	if len(left) != 0 {
		return nil, errors.New("trailing input " + string(left))
	}

	return result, nil
}

func parseExpr(input []byte) (result *expr, left []byte, err error) {
	_, input, _ = trimWhitespace(input)
	if input[0] == '(' {
		return parseOperation(input)
	} else if isDigit(input[0]) {
		return parseLeaf(input)
	}
	return nil, input, errors.New("expected a number or an arithmetic expression")
}

func parseOperation(in []byte) (*expr, []byte, error) {
	if in[0] != '(' {
		return nil, in, errors.New("expected opening parenthesis")
	}

	_, in, _ = trimWhitespace(in[1:])
	var operator Operation
	var present bool
	if operator, present = operationForChar[in[0]]; !present {
		return nil, in, errors.New("unknown operator")
	}
	in = in[1:]

	var operands []expr = make([]expr, 0, 10)
	var operand *expr
	var err error
	for len(in) > 0 && in[0] != ')' {
		operand, in, err = parseExpr(in)
		if err != nil {
			return nil, in, err
		}
		operands = append(operands, *operand)
		_, in, _ = trimWhitespace(in)
	}
	if len(in) == 0 || in[0] != ')' {
		return nil, in, errors.New("expected closing parenthesis")
	}
	return Expression(operator, operands), in[1:], nil
}

func parseLeaf(in []byte) (*expr, []byte, error) {
	var i int
	for i = 0; i < len(in); i++ {
		if !isDigit(in[i]) {
			break
		}
	}
	leafValue, err := strconv.ParseInt(string(in[0:i]), 10, 64)
	if err != nil {
		return nil, in, err
	}
	return Leaf(leafValue), in[i:], nil
}

var AllOperations []Operation

func init() {
	operationForChar = make(map[byte]Operation)
	AllOperations = []Operation{
		Addition{},
		Subtraction{},
		Multiplication{},
		Division{}}
	for _, op := range AllOperations {
		operationForChar[byte(op.String()[0])] = op
	}
}
