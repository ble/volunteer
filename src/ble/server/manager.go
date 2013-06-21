package server

import (
  "container/ring"
  "errors"
)

type Manager interface {
  volunteerAdder(Worker)
  nextAdder() (Worker, error)
  dropAdder(Worker)

  volunteerSubber(Worker)
  nextSubber() (Worker, error)
  dropSubber(Worker)

  volunteerMuller(Worker)
  nextMuller() (Worker, error)
  dropMuller(Worker)

  volunteerDivver(Worker)
  nextDivver() (Worker, error)
  dropDivver(Worker)
}

type managerImpl struct {
  
}


