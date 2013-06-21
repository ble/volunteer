package main

import (
	"ble/commander"
	"ble/server"
	"log"
	"net/http"
)

func main() {
	manager := server.SetUpServer()
	commander.ConfigureHandlers(manager)
	err := http.ListenAndServe(":24736", nil)
	log.Println(err)
}
