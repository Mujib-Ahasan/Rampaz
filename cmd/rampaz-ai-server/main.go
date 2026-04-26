package main

import (
	"log"
	"net/http"

	"github.com/Mujib-Ahasan/Rampaz/internal/ai"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	server := ai.NewServer()
	defer server.Close()

	log.Println("Rampaz AI server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", server.Routes()); err != nil {
		log.Fatal(err)
	}
}
