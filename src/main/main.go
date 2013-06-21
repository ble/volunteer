package main

import (
  "ble/server"
  "log"
  "net/http"
)

func main() {
  _ = server.SetUpServer()
  err := http.ListenAndServe(":24736", nil)
  log.Println(err)
}
