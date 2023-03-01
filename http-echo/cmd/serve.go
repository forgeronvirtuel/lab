package main

import (
	"flag"
	"io"
	"log"
	"net/http"
)

func main() {
	// A simple command line to precise the host
	host := flag.String("host", "localhost:8080", "listening address and port")
	flag.Parse()

	// The main route repeat the content of the request
	handler := func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(w, r.Body); err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
	}
	http.HandleFunc("/", handler)

	log.Printf("Listening to : %s", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal(err)
	}
}
