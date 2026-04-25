package main

import (
	"log"
	"net/http"

	"github.com/Mujib-Ahasan/Rampaz/internal/ai"
)

func main() {
	server := ai.NewServer()

	log.Println("Rampaz AI server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", server.Routes()); err != nil {
		log.Fatal(err)
	}
}
