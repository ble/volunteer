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
		h.ServeFile(w, r, "static/commandClient.html")
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

//kind of lame way of doing it; ideally we'd omit the value field when
//Error != ""
type Completion struct {
	Problem string `json:"problem"`
	Value   int64  `json:"value"`
	Error   string `json:"error,omitempty"`
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
			response := c.evaluate(*parsed)
			var eString string
			if response.Error != nil {
				eString = response.Error.Error()
			}
			encoder.Encode(Completion{s.Problem, response.Value, eString})
		}

	}
}
