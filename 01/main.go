package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var isHealth bool

func main() {
	isHealth = true
	http.HandleFunc("/health", healthyHandler)
	http.HandleFunc("/unhealth", unhealthyHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/", indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
		log.Printf("Listening on port %s", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	}
}

func healthyHandler(w http.ResponseWriter, r *http.Request) {
	isHealth = true
	_, err := fmt.Fprintf(w, "change response code to 200")
	if err != nil {
		w.WriteHeader(http.StatusOK)
	}
}

func unhealthyHandler(w http.ResponseWriter, r *http.Request) {
	isHealth = false
	_, err := fmt.Fprintf(w, "change response code to 503")
	if err != nil {
		w.WriteHeader(http.StatusOK)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if isHealth == true {
		_, err := fmt.Fprintf(w, "Healty!!!")
		if err != nil {
			w.WriteHeader(http.StatusOK)
		}
		return
	}

	w.WriteHeader(http.StatusServiceUnavailable)
	return
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "pong")
	if err != nil {
		w.WriteHeader(http.StatusOK)
	}
}
