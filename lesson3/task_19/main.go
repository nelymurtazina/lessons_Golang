package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type RequestData struct {
	UserId int      `json:"user_id"`
	Action string   `json:"action"`
	Items  []string `json:"items"`
}

func (r *RequestData) Reset() {
	r.UserId = 0
	r.Action = ""
	// Сохраняем базовый массив, обнуляем длину
	r.Items = r.Items[:0]
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return &RequestData{
			Items: make([]string, 0, 10), //аллокация
		}
	},
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// закрыть body
	defer r.Body.Close()

	// Получаем объект из пула
	data := requestPool.Get().(*RequestData)
	defer requestPool.Put(data)

	data.Reset()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(data); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return 
	}

	fmt.Printf("Обработан запрос: UserId=%d, Action=%s, Items=%v\n",
		data.UserId, data.Action, data.Items)

	// ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		fmt.Printf("Ошибка отправки ответа: %v\n", err)
	}
}

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Сервер остановлен: %v\n", err)
	}
}