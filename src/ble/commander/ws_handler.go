package commander

import (
  h "net/http"
  "ble/parse"
  "ble/server"
  "code.google.com/p/go.net/websocket"
  "encoding/json"
)

type Evaluator interface{}
type Commander struct {server.Manager}

func configureHandlers(m server.Manager) Commander {
  c := Commander{m}
  h.Handle("/command", websocket.Handler(c.wsHandler()))
  h.HandleFunc("/commandClient", func(w h.ResponseWriter, r *h.Request) {
    w.Write([]byte(
`
<html><head></head><body>
<form id="inputForm">
  <input type="text"></input>
  <input type="submit"></input>
</form>

<script>
  var url = location.toString();
  var wsUrl = url.replace(/^http/,"ws").replace("Client","")
  var socket = new WebSocket(wsUrl);

  var inputForm = document.getElementById('inputForm');
  inputForm.addEventListener('submit', function(e) {
    window.console.log(e);
    e.preventDefault();
    window.console.log(inputForm);
      
  });

</script>

</body></html>
`))
  })
  return c
}

type Problem struct {
  PolishNotation string `json:"polishNotation"`
}

type Status struct {
  Status string `json:"status"`
}

func (c Commander) wsHandler() websocket.Handler {
  return func(conn *websocket.Conn) {
    defer conn.Close()
    decoder := json.NewDecoder(conn)
    encoder := json.NewEncoder(conn)

    p := Problem{}
    for {
      err := decoder.Decode(p)
      if err != nil {
        encoder.Encode(Status{"json error " + err.Error()})
        break
      } else {
        encoder.Encode(Status{"received"})
      }

      parsed, err := parse.Parse([]byte(p.PolishNotation))
      if err != nil {
        encoder.Encode(Status{"polish notation error " + err.Error()})
        break
      } else {
        encoder.Encode(Status{"parsed"})
      }
      println(parsed)
    }

  }
}
