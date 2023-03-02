package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// A simple command line to precise the host
	host := flag.String("url", "http://localhost:8080/", "Address and port to call")
	content := flag.String("content", "Hello, World!", "Content to send to the url through a POST request.")
	flag.Parse()

	requestBody := []byte(*content)

	// Send the content to the backend server
	resp, err := http.Post(*host, "text/plain", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Error sending request: %s\n", err.Error())
		return
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
