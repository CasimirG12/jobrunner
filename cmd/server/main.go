package main

import (
	"log"
	"net/http"

	"github.com/CasimirG12/jobrunner/internal/api"
	"github.com/CasimirG12/jobrunner/internal/store"
)

type ResponseMessage struct {
	Message string `json:"message"`
}

func main() {
	store := store.NewMemoryStore()
	server := api.NewServer(store)
	routes := server.Routes()

	log.Println("Listening to port 8080...")
	log.Fatal(http.ListenAndServe(":8080", routes))
}
