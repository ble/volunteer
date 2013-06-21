package volunteer

import (
  h "net/http"
  . "ble/parse"
)

func ConfigureWSHandlers(m Manager) {
  done := make(chan Worker)
  table := []struct{path string; Operation} {
    {"/volunteer/add", AllOperations[0]},
    {"/volunteer/sub", AllOperations[1]},
    {"/volunteer/mul", AllOperations[2]},
    {"/volunteer/div", AllOperations[3]},
  }

  for _, tableEntry := range table {
    h.HandleFunc(tableEntry.path, func(w h.ResponseWriter, r *h.Request) {
      wsHandler, worker := MakeWorkerHandler(tableEntry.Operation, done)
      m.volunteer(tableEntry.Operation, worker)
      wsHandler.ServeHTTP(w, r)
    })
  }

}
