package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// A simple command line to precise the host
	host := flag.String("host", "localhost:8080", "listening address and port")
	flag.Parse()

	incrementCountCh := make(chan struct{}, 10)
	askForCountCh := make(chan struct{})
	sendCountCh := make(chan int)

	// Handle the counts
	go func() {
		var count int
		for {
			select {
			case <-incrementCountCh:
				count += 1
			case <-askForCountCh:
				sendCountCh <- count
			}
		}
	}()

	// The main route repeat the content of the request
	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		incrementCountCh <- struct{}{}
		if _, err := io.Copy(w, r.Body); err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
	}

	// Count handler send the number of query called
	countHandler := func(w http.ResponseWriter, r *http.Request) {
		askForCountCh <- struct{}{}
		count := <-sendCountCh
		if _, err := w.Write([]byte(fmt.Sprintf("%d", count))); err != nil {
			http.Error(w, "Error writing response body", http.StatusInternalServerError)
		}
	}

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/count/", countHandler)

	log.Printf("Listening to : %s", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal(err)
	}
}
