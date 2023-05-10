package net

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type OutgoingEventStream struct {
	// Request made to initiate the connection with the server.
	ConnectionRequest http.Request
}

// Create a new `OutgoingEventStream`, please denote that `ConnectionOpen` should be checked often to know if the connection is still open
// or not. To create this instance of the struct, it operates the same as a classic HTTP binder but it opens instead a SSE connection.
func MakeOutgoingEventStreamFromRequest(req *http.Request, appropriateBinding *Binding) (*OutgoingEventStream, error) {
	// Extract the port to target
	var port, err = ExtractBindedPort(req)
	if err != nil {
		return nil, errors.New("no " + TARGET_HEADER_NAME + " header found")
	}

	// Build the request
	req.RequestURI = ""
	req.URL.Host = fmt.Sprintf("localhost:%d", port)

	// Determines which scheme to use due to the scheme override setting.
	if appropriateBinding.DefaultTransfertScheme != nil {
		req.URL.Scheme = *appropriateBinding.DefaultTransfertScheme
	} else {
		req.URL.Scheme = "http"
	}

	req.Header.Set("Accept", "text/event-stream")

	return &OutgoingEventStream{ConnectionRequest: *req}, nil
}

// Create a new `OutgoingEventStream`, please denote that `ConnectionOpen` should be checked often to know if the connection is still open
// or not.
func MakeOutgoingEventStreamFromURL(url string) (*OutgoingEventStream, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "text/event-stream")

	return &OutgoingEventStream{ConnectionRequest: *req}, nil
}

// Make and read for events on the built connection
func (oes *OutgoingEventStream) ReadEvents(listener func(readLine string)) error {
	// Sends the request
	res, err := http.DefaultClient.Do(&oes.ConnectionRequest)
	if err != nil {
		return errors.New("network error")
	}

	// Make appropriate readers for the request
	rbuf := bufio.NewReader(res.Body)
	defer res.Body.Close()

	for {
		revent, err := rbuf.ReadBytes('\n')

		// Fails if the read line doesn't meet requirements
		if err != nil && err != io.EOF {
			return err
		}

		// Give the line to the listener that has to handle it.
		listener(string(revent))
	}
}
