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
	http.HandleFunc("/volunteerClient", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(
			`<html><body>
<script>
var url = window.location.toString();
var wsUrl = url.replace(/^http/, "ws").replace("volunteerClient", "volunteer/add");
var ws = new WebSocket(wsUrl);
</script>
</body></html>`))
	})
	http.Handle("/echo", websocket.Handler(EchoHandler))
	http.HandleFunc("/echoClient", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(
			`<html><body>
<script>
var url = window.location.toString();
var wsUrl = url.replace(/^http/, "ws").replace("clientPage", "echo");
var ws = new WebSocket(wsUrl);
ws.addEventListener('open', function(e) {
  var socket = e.target;
  socket.messagesSent = 0;
  socket.send("hello");
  socket.messagesSent++;
});
ws.addEventListener('message', function(e) {
  var socket = e.target;
  window.console.log(e);
  socket.send("ping");
  socket.messagesSent++;
  if(socket.messagesSent > 5)
    socket.close();
});
ws.addEventListener('close', function(e) {
  var socket = e.target;
  window.console.log(e);
  window.console.log("closing");
});
</script>
</body></html>`))
	})
	return m
}

func EchoHandler(conn *websocket.Conn) {
	_, err := io.Copy(conn, conn)
	if err != nil {
		log.Println("Echo handler: " + err.Error())
	}
}
