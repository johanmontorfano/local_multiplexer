package net

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Used to manage incoming event streams (receiving data); and outgoing event
// streams (sending data).
type EventStream struct {
	SupportEventStreams bool
	flusher             *http.Flusher
	ResponseWriter      *http.ResponseWriter
    InitializerReq      *http.Request
}

// Creates a new `EventStream`, please denote that the only way to know if it 
// has been created successfully is by reading the `SupportEventStreams` field.
// Informing the client of a possible error should be handled manually.
func MakeEventStream(
    w http.ResponseWriter, 
    req *http.Request, 
    bind Binding,
) EventStream {
    req.RequestURI = ""
    req.URL.Host = fmt.Sprintf("localhost:%d", bind.TargetPort)
    req.Header.Set("Accept", "text/event-stream")
    req.URL.Scheme = "http"

    // Determines if the default transfert scheme should be overwritten.
    if bind.TransfertProtocol != nil {
        req.URL.Scheme = *bind.TransfertProtocol
    }

    flusher, ok := w.(http.Flusher)

	if ok {
		w.Header().Add("Content-Type", "text/event-stream")
		w.Header().Add("Connection", "keep-alive")
		w.WriteHeader(200)
		flusher.Flush()

		return EventStream {
            SupportEventStreams: ok, 
            flusher: &flusher, 
            ResponseWriter: &w,
            InitializerReq: req,
        }
	} else {
		return EventStream {
            SupportEventStreams: ok, 
            flusher: nil, 
            ResponseWriter: nil,
            InitializerReq: nil,
        }
	}
}

// Listen for events sent through SSE by the client. May do nothing if 
// `SupportEventStreams` is false.
func (es *EventStream) ReadEvents (callback func(line string)) error {
    res, err := http.DefaultClient.Do(es.InitializerReq)

    if err != nil {
        return errors.New("network error")
    }

    readerbuf := bufio.NewReader(res.Body)
    defer res.Body.Close()

    for {
        event, err := readerbuf.ReadBytes('\n')

        if err != nil && err != io.EOF { return err }

        callback(string(event))
    }
}

// Add data to the next buffer to flush. May do nothing if `SupportEventStreams` 
// is false. To fetch the written data to the client, please use the 
// `IncomingEventStream.SendBuffer` function.
func (es *EventStream) Add2Buffer(data string) {
	if es.ResponseWriter != nil {
		(*es.ResponseWriter).Write([]byte(data))
	}
}

// Sends the buffered data to the client, the function may do nothing if 
// `SupportEventStreams` is not took in considerations.
func (es *EventStream) SendBuffer() {
	if es.flusher != nil {
		(*es.flusher).Flush()
	}
}
