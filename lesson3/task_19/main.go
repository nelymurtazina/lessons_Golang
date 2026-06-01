package main

import (
	"sync"
	"net/http"
	"fmt"
	"encoding/json"
)

type RequestData struct {
	UserId int
	Action string
	Items  []string
}

func (r *RequestData) Reset() {
	r.UserId = 0
	r.Action = ""
	r.Items = r.Items[:0]
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return &RequestData{
			Items: make([]string, 0, 10),
		}
	},
}

func handleRequest (w http.ResponseWriter, r *http.Request) {
	data := requestPool.Get().(*RequestData)
	defer requestPool.Put(data)

	data.Reset()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(data); err !=nil {
		http.Error(w,"Bad request", http.StatusAccepted)
	}

	fmt.Println("Обработан запрос", data.UserId, data.Action, data.Items)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}