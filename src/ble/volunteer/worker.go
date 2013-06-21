package volunteer


import (
  "ble/parse"
  "code.google.com/p/go.net/websocket"
  "encoding/json"
  "log"
)

type Worker interface {
  Evaluate(e parse.Expr) (int64, error)
  Operator() parse.Operation
}

type workerImpl struct {
  parse.Operation
  *websocket.Conn
  requests chan request
  done chan<- Worker
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
  expr parse.Expr
  result chan Response
}

type Response struct {
  Value int64 `json:"value"`
  Error error `json:"-"`
}

func MakeWorkerHandler(o parse.Operation, done chan<- Worker) (websocket.Handler, Worker) {
  ch := make(chan request)
  w := &workerImpl{o, nil, ch, done}
  h := func(conn *websocket.Conn) {
    defer conn.Close()
    encoder := json.NewEncoder(conn)
    decoder := json.NewDecoder(conn)

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

      err = decoder.Decode(resp)
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
