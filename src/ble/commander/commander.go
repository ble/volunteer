package commander

import (
  "ble/volunteer"
  "ble/parse"
)

type EvaluationRequest struct {
  parse.Expr
  reply chan volunteer.Response
}

type Commander struct {
  volunteer.Manager
  requests chan EvaluationRequest
}
