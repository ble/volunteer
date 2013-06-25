package volunteer

import (
	. "ble/parse"
	h "net/http"
)

func ConfigureWSHandlers(m Manager) {
	done := make(chan Worker)
	table := []struct {
		path string
		Operation
	}{
		{"/volunteer/add", AllOperations[0]},
		{"/volunteer/sub", AllOperations[1]},
		{"/volunteer/mul", AllOperations[2]},
		{"/volunteer/div", AllOperations[3]},
	}

	for _, tableEntry := range table {
		//not making a copy === classic mistake in creating closures inside of a
		//loop.
		tableEntryCopy := tableEntry
		h.HandleFunc(tableEntry.path, func(w h.ResponseWriter, r *h.Request) {
			wsHandler, worker := MakeWorkerHandler(tableEntry.Operation, done)
			m.volunteer(tableEntryCopy.Operation, worker)
			wsHandler.ServeHTTP(w, r)
		})
	}
	h.HandleFunc("/volunteerClient", func(w h.ResponseWriter, r *h.Request) {
		h.ServeFile(w, r, "static/volunteerClient.html")
	})

}
