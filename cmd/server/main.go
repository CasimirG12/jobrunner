package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseMessage struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, v any, code int) {
	body, err := json.Marshal(v)

	if err != nil {
		log.Printf("marshal response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, ResponseMessage{Message: "Hello from the server!"}, http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
