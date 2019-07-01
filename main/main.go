package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"../config"
	"../handlers"

	"src/github.com/gorilla/mux"
)

func main() {
	os.Setenv("PORT", "8081")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	config := config.New("conf.yaml")

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HandleMain)
	r.HandleFunc("/products", handlers.HandleProduct)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
