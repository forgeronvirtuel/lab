package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// A simple command line to precise the host
	host := flag.String("url", "http://localhost:8080/", "Address and port to call")
	flag.Parse()

	// Send the content to the backend server
	countUrl, err := url.JoinPath(*host, "count/")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(countUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err.Error())
		return
	}

	// Print the response body.
	fmt.Println(string(body))
}
