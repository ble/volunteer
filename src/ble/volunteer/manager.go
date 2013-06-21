package volunteer

import (
  "ble/parse"
  "container/ring"
  "errors"
  "sync"
)

type Manager interface {
  volunteer(parse.Operation, Worker)
  next(parse.Operation) (Worker, error)
  drop(parse.Operation, Worker)
}

type managerImpl struct {
  workersByOp map[parse.Operation]*ring.Ring
  sync.Mutex
}

func NewManager() Manager {
  return &managerImpl{make(map[parse.Operation]*ring.Ring), sync.Mutex{}}
}

func (m *managerImpl) volunteer(o parse.Operation, w Worker) {
  println("New volunteer!")
  m.Lock()
  defer m.Unlock()

  //since a nil *Ring is an empty ring, no need to check for presence
  currentWorkers := m.workersByOp[o]

  //make a link for the new worker, attach it before current workers
  newWorker := ring.New(1)
  newWorker.Value = w
  newWorker.Link(currentWorkers)

  //if currentWorkers wasn't empty, its current element will be the same
  m.workersByOp[o] = newWorker.Next()
}

func (m *managerImpl) next(o parse.Operation) (Worker, error) {
  m.Lock()
  defer m.Unlock()

  wRing := m.workersByOp[o]
  if wRing.Len() == 0 {
    return nil, errors.New("no workers of requested type")
  }

  worker, ok := wRing.Value.(Worker)
  if !ok {
    return nil, errors.New("someone snuck a non-worker into our group...")
  }

  return worker, nil
}

func (m *managerImpl) drop(o parse.Operation, w Worker) {
  m.Lock()
  defer m.Unlock()

  targetRing := m.workersByOp[o]

  for i := 0; i < targetRing.Len(); i++ {
    w_i, _ := targetRing.Value.(Worker)
    if w_i == w {
      targetRing.Prev().Unlink(1)
      break
    }
    targetRing = targetRing.Next()
  }
}
