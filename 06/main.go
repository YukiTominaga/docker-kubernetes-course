package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var isHealth bool
var hostname = os.Getenv("HOSTNAME")

func main() {
	if os.Getenv("HEALTHY") == "FALSE" || os.Getenv("HEALTHY") == "false" {
		isHealth = false
	} else {
		isHealth = true
	}

	http.HandleFunc("/health", healthyHandler)
	http.HandleFunc("/unhealth", unhealthyHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/", indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
		log.Printf("HostName is %s", hostname)
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func healthyHandler(w http.ResponseWriter, r *http.Request) {
	isHealth = true
	_, err := fmt.Fprintf(w, "%s: change response code to 200", hostname)
	if err != nil {
		w.WriteHeader(http.StatusOK)
	}
}

func unhealthyHandler(w http.ResponseWriter, r *http.Request) {
	isHealth = false
	_, err := fmt.Fprintf(w, "%s: change response code to 503", hostname)
	if err != nil {
		w.WriteHeader(http.StatusOK)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if isHealth == true {
		_, err := fmt.Fprintf(w, "HostName is %s", hostname)
		if err != nil {
			w.WriteHeader(http.StatusOK)
		}
		return
	}

	w.WriteHeader(http.StatusServiceUnavailable)
	return
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "%s: pong", hostname)
	if err != nil {
		w.WriteHeader(http.StatusOK)
	}
}
