package commander

import (
	"ble/parse"
	"ble/volunteer"
)

type EvaluationRequest struct {
	parse.Expr
	reply chan volunteer.Response
}

type Commander struct {
	volunteer.Manager
	requests chan EvaluationRequest
}
