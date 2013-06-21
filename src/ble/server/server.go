package server

import (
	"ble/volunteer"
	"code.google.com/p/go.net/websocket"
	"io"
	"log"
	"net/http"
)

func SetUpServer() volunteer.Manager {
	m := volunteer.NewManager()

	volunteer.ConfigureWSHandlers(m)
	http.Handle("/echo", websocket.Handler(EchoHandler))
	http.HandleFunc("/echoClient", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/echoClient.html")
	})
	return m
}

func EchoHandler(conn *websocket.Conn) {
	_, err := io.Copy(conn, conn)
	if err != nil {
		log.Println("Echo handler: " + err.Error())
	}
}
