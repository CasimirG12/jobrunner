package main

import (
	"log"
	"net/http"

	"github.com/CasimirG12/jobrunner/internal/api"
)

type ResponseMessage struct {
	Message string `json:"message"`
}

func main() {
	server := api.NewServer()
	routes := server.Routes()

	log.Fatal(http.ListenAndServe(":8080", routes))
}
