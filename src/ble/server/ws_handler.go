package server

import (
  h "net/http"
  . "ble/parse"
)

func configureWSHandlers(m Manager) {
  done := make(chan Worker)
  table := []struct{path string; Operation; volunteer func(m Manager, w Worker)} {
    {"/volunteer/add", Addition{}, (Manager).volunteerAdder},
    {"/volunteer/sub", Subtraction{}, (Manager).volunteerSubber},
    {"/volunteer/mul", Multiplication{}, (Manager).volunteerMuller},
    {"/volunteer/div", Division{}, (Manager).volunteerDivver},
  }

  for _, volunteerType := range table {
    h.HandleFunc(volunteerType.path, func(w h.ResponseWriter, r *h.Request) {
      wsHandler, worker := MakeWorkerHandler(volunteerType.Operation, done)
      wsHandler.ServeHTTP(w, r)
      volunteerType.volunteer(m, worker)
    })
  }

}
