package commander

import (
	"ble/parse"
	"ble/volunteer"
	"errors"
	"fmt"
)

type EvaluationRequest struct {
	parse.Expr
	reply chan volunteer.Response
}

type Commander struct {
	volunteer.Manager
	requests chan EvaluationRequest
}

func (c Commander) evaluate(topLevel parse.Expr) volunteer.Response {
	opWorkers := []volunteer.Worker{nil, nil, nil, nil}
	var err error

	//define a type so that we can manage our own stack rather than write this
	//as a recursive function.
	type stackFrame struct {
		current, parent *parse.Expr
		operandIndex    int
	}
	stack := []stackFrame{{&topLevel, nil, -1}}

	for {
		frame := stack[len(stack)-1]
		fmt.Printf("%#v\n", frame)
		operation := frame.current.Operation

		//ensure that we have a worker for this particular operation
		workerIndex := parse.OperationIndex(operation)
		if workerIndex == -1 {
			return volunteer.Response{0, errors.New("couldn't match operation")}
		}
		if opWorkers[workerIndex] == nil {
			//hold on to this worker for all operations of this type
			opWorkers[workerIndex], err = c.Manager.Next(operation)
			if err != nil {
				return volunteer.Response{0, errors.New("Manager.Next: " + err.Error())}
			}
		}
		worker := opWorkers[workerIndex]

		//send the expression to the worker if it has no child subexpressions
		if frame.current.NoGrandChildren() {
			println("sending to volunteer...")
			value, err := worker.Evaluate(*frame.current)

			//evaluation failed: we give up on the whole thing
			if err != nil {
				println("error on worker.Evaluate")
				return volunteer.Response{0, err}
			}

			//we just evaluated a sub-expression and need to substitute the reuslt
			//into the containing expression
			if frame.parent != nil {
				frame.parent.Operands[frame.operandIndex] = parse.Leaf(value)

				//pop dat stack
				stack = stack[:len(stack)-1]

				//are we there yet?
				if len(stack) == 0 {
					return volunteer.Response{value, nil}
				}
			}
			println("evaluated!")
		} else {
			println("pushing onto stack")
			//we have 1+ child subexpressions and need to evaluate them first
			for index, subExpression := range frame.current.Operands {

				//if it's a subexpression, it goes on the stack.
				if !subExpression.IsLeaf() {
					stack = append(stack, stackFrame{&subExpression, frame.current, index})
				}
			}
		}
	}
	return volunteer.Response{0, errors.New("didn't feel like evaluating")}
}
