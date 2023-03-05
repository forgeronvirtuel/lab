package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// A simple command line to precise the host
	host := flag.String("host", "localhost:8080", "listening address and port")
	flag.Parse()

	srv := &http.Server{Addr: *host}

	wg := sync.WaitGroup{}

	// Send a signal to increment the internal counter
	incrementCountCh := make(chan struct{}, 10)
	// Send a signal to send the current counter value
	askForCountCh := make(chan struct{})
	// The channel where we send the current value of the counter
	sendCountCh := make(chan int)
	// Capture the SIGINT signal to gracefully stop the server and all goroutine
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	// Handle the counts and the closing server event
	wg.Add(1)
	go func() {
		var count int
		var stop bool
		for !stop {
			select {
			case <-incrementCountCh:
				count += 1
			case <-askForCountCh:
				sendCountCh <- count
			case <-sigCh:
				stop = true
				if err := srv.Close(); err != nil {
					log.Fatal(err)
				}
			}
		}
		wg.Done()
	}()

	// The main route repeat the content of the request
	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		incrementCountCh <- struct{}{}
		if _, err := io.Copy(w, r.Body); err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
	}
	http.HandleFunc("/", mainHandler)

	// Count handler send the number of query called
	countHandler := func(w http.ResponseWriter, r *http.Request) {
		askForCountCh <- struct{}{}
		count := <-sendCountCh
		if _, err := w.Write([]byte(fmt.Sprintf("%d", count))); err != nil {
			http.Error(w, "Error writing response body", http.StatusInternalServerError)
		}
	}
	http.HandleFunc("/count/", countHandler)

	log.Printf("Listening to : %s", *host)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	wg.Wait()
}
