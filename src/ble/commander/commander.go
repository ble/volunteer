package commander

import (
	"ble/parse"
	"ble/volunteer"
	"errors"
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

	type stackFrame struct {
		current, parent *parse.Expr
		operandIndex    int
	}
	stack := []stackFrame{{&topLevel, nil, -1}}

	//return volunteer.Response{0, errors.New("didn't feel like evaluating")}
	for len(stack) > 0 {
		frame := stack[len(stack)-1]
		operation := frame.current.Operation

		//ensure that we have a worker for this particular operation
		workerIndex := parse.OperationIndex(operation)
		if workerIndex == -1 {
			return volunteer.Response{0, errors.New("couldn't match operation")}
		}
		if opWorkers[workerIndex] == nil {
			opWorkers[workerIndex], err = c.Manager.Next(operation)
			if err != nil {
				return volunteer.Response{0, errors.New("Manager.Next: " + err.Error())}
			}
		}
		worker := opWorkers[workerIndex]

		//send the expression to the worker if it has no child subexpressions
		if frame.current.NoGrandChildren() {
			value, err := worker.Evaluate(*frame.current)
			return volunteer.Response{value, err}
		} else {
			//add all the subexpressions
		}
		return volunteer.Response{0, errors.New("didn't feel like evaluating")}
	}
	return volunteer.Response{0, errors.New("didn't feel like evaluating")}
}
