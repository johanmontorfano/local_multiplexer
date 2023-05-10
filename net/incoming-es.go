package net

import "net/http"

type IncomingEventStream struct {
	SupportEventStreams bool
	flusher             *http.Flusher
	ResponseWriter      *http.ResponseWriter
}

// Creates a new `IncomingEventStream`, please denote that the only way to know if it has been created successfully is by
// reading the `SupportEventStreams` field. Informing the client of a possible error should be handled manually.
func MakeIncomingEventStream(w http.ResponseWriter) IncomingEventStream {
	flusher, ok := w.(http.Flusher)

	if ok {
		w.Header().Add("content-type", "text/event-stream")
		w.Header().Add("connection", "keep-alive")
		w.WriteHeader(200)
		flusher.Flush()

		return IncomingEventStream{SupportEventStreams: ok, flusher: &flusher, ResponseWriter: &w}
	} else {
		return IncomingEventStream{SupportEventStreams: ok, flusher: nil, ResponseWriter: nil}
	}
}

// Writes data to a connected client, the function may do nothing if `SupportEventStreams` is not took in considerations.
// To fetch the written data to the client, please use the `IncomingEventStream.SendBuffer` function.
func (ies *IncomingEventStream) WriteToBuffer(data string) {
	if ies.ResponseWriter != nil {
		(*ies.ResponseWriter).Write([]byte(data))
	}
}

// Sends the buffered data to the client, the function may do nothing if `SupportEventStreams` is not took in considerations.
func (ies *IncomingEventStream) SendBuffer() {
	if ies.flusher != nil {
		(*ies.flusher).Flush()
	}
}
