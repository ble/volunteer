package commander

import (
	"ble/parse"
	"ble/volunteer"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"log"
	h "net/http"
)

func ConfigureHandlers(m volunteer.Manager) Commander {
	c := Commander{m, make(chan EvaluationRequest)}
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

  socket.addEventListener('message', function(e) {
    window.console.log(e.data);
    var el = document.createElement("div");
    el.appendChild(document.createTextNode(e.data));
    document.body.appendChild(el);
  });

  var inputForm = document.getElementById('inputForm');
  inputForm.addEventListener('submit', function(e) {
    var inputText = e.target.elements[0];
    window.console.log(e);
    e.preventDefault();
    var jsonObj = {'polishNotation': inputText.value};
    socket.send(JSON.stringify(jsonObj));
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
	Problem string `json:"problem"`
	Status  string `json:"status"`
}

func (c Commander) wsHandler() websocket.Handler {
	return func(conn *websocket.Conn) {
		println("awaiting commands")
		defer conn.Close()
		defer println("command connection closed.")
		decoder := json.NewDecoder(conn)
		encoder := json.NewEncoder(conn)

		p := Problem{}
		s := Status{"", ""}
		for {
			err := decoder.Decode(&p)
			log.Println("decoded!")
			if err != nil {
				encoder.Encode(Status{"json error " + err.Error(), ""})
				break
			} else {
				s.Problem = p.PolishNotation
				s.Status = "received"
				encoder.Encode(s)
			}

			parsed, err := parse.Parse([]byte(p.PolishNotation))
			log.Println("parsed!")
			if err != nil {
				s.Status = "polish notation error " + err.Error()
				encoder.Encode(s)
				break
			} else {
				s.Status = "parsed"
				encoder.Encode(s)
			}
			println(parsed.String())
		}

	}
}
