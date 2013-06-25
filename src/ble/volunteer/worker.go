package volunteer

import (
	"ble/parse"
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
)

type Worker interface {
	Evaluate(e parse.Expr) (int64, error)
	Operator() parse.Operation
}

type workerImpl struct {
	parse.Operation
	*websocket.Conn
	requests chan request
	done     chan<- Worker
}

func (w *workerImpl) Evaluate(e parse.Expr) (int64, error) {
	replyChannel := make(chan Response)
	req := request{e, replyChannel}
	w.requests <- req
	reply := <-replyChannel
	return reply.Value, reply.Error
}

func (w *workerImpl) Operator() parse.Operation {
	return w.Operation
}

type request struct {
	expr   parse.Expr
	result chan Response
}

type Response struct {
	Value int64 `json:"value"`
	Error error `json:"error,omitempty"`
}

func (r *Response) MarshalJSON() ([]byte, error) {
	substitute := make(map[string]interface{})
	substitute["value"] = r.Value
	if r.Error != nil {
		substitute["error"] = r.Error.Error()
	}
	return json.Marshal(substitute)
}

func (r *Response) UnmarshalJSON(data []byte) error {
	substitute := make(map[string]interface{})
	if err := json.Unmarshal(data, &substitute); err != nil {
		return err
	}
	if value, present := substitute["value"]; present {
		//this is monumentally stupid: it will panic if the key "value" has a non-
		//int value.
		r.Value = reflect.ValueOf(value).Int()
	} else {
		return errors.New("no value present in response")
	}
	if errString, present := substitute["error"]; present {
		r.Error = errors.New(reflect.ValueOf(errString).String())
	} else {
		r.Error = nil
	}
	return nil
}

func MakeWorkerHandler(o parse.Operation, done chan<- Worker) (websocket.Handler, Worker) {
	ch := make(chan request)
	w := &workerImpl{o, nil, ch, done}

	//awkward as hell, for debugging.
	jsonBytes := make([]byte, 1024, 1024)
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))

	h := func(conn *websocket.Conn) {
		defer conn.Close()
		encoder := json.NewEncoder(conn)
		decoder := json.NewDecoder(buffer)

		for {
			var resp Response

			req := <-w.requests
			err := encoder.Encode(req.expr)
			if err != nil {
				log.Println("sending json: " + err.Error())
				resp.Error = err
				req.result <- resp
				break
			}

			n, err := conn.Read(jsonBytes)
			buffer.Write(jsonBytes[0:n])
			fmt.Printf("read data: '%s'", string(jsonBytes))
			err = decoder.Decode(&resp)
			fmt.Printf("worker response: %#v\n", resp)
			if err != nil {
				log.Println("receiving json: " + err.Error())
				resp.Error = err
				req.result <- resp
				break
			}
			req.result <- resp
		}
		done <- w
	}
	return h, w
}
