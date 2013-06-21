package parse

import (
	"encoding/json"
	"fmt"
  "strings"
)

type Operation interface {
	Evaluate(operands []int64) int64
	String() string
}

type Expr interface {
  Operator() Operation
  IsLeaf() bool
  NoGrandChildren() bool
  fmt.Stringer
}

type expr struct {
	Operands []expr
	Operation
	LeafValue int64
}

func (e expr) Operator() Operation {
  return e.Operation
}

func (e expr) IsLeaf() bool {
	return e.Operation == nil
}

func (e expr) NoGrandChildren() bool {
  for _, child := range(e.Operands) {
    if !child.IsLeaf() {
      return false
    }
  }
  return true
}

func (e expr) String() string {
  if(e.IsLeaf()) {
    return fmt.Sprintf("%d", e.LeafValue)
  }
  parts := make([]string, 0, 2+len(e.Operands))
  parts = append(parts, "(" + e.Operation.String())
  for _, operand := range e.Operands {
    parts = append(parts, operand.String())
  }
  parts = append(parts, ")")
  return strings.Join(parts, " ")
}

func Leaf(x int64) *expr {
	return &expr{nil, nil, x}
}

func Expression(o Operation, os []expr) *expr {
	return &expr{os, o, 0}
}

func (e expr) MarshalJSON() ([]byte, error) {
	if e.IsLeaf() {
		return []byte(fmt.Sprintf("%d", e.LeafValue)), nil
	}
	tmp := make(map[string]interface{})
	tmp["operator"] = e.Operation.String()
	tmp["operands"] = e.Operands
	return json.Marshal(tmp)
}

type Addition struct{}

func (a Addition) Evaluate(operands []int64) int64 {
	var result int64 = 0
	for _, operand := range operands {
		result += operand
	}
	return result
}
func (a Addition) String() string { return "+" }

type Subtraction struct{}

func (s Subtraction) Evaluate(operands []int64) int64 {
	var result int64
	if len(operands) != 1 {
		result = operands[0]
		for _, operand := range operands[1:] {
			result -= operand
		}
	} else {
		result = -operands[0]
	}
	return result
}
func (s Subtraction) String() string { return "-" }

type Multiplication struct{}

func (a Multiplication) Evaluate(operands []int64) int64 {
	var result int64 = 1
	for _, operand := range operands {
		result *= operand
	}
	return result
}
func (m Multiplication) String() string { return "*" }

type Division struct{}

func (s Division) Evaluate(operands []int64) int64 {
	result := operands[0]
	for _, operand := range operands[1:] {
		result -= operand
	}
	return result
}
func (d Division) String() string { return "/" }
