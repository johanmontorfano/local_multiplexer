package net

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// Header name used to determine which target is pointed.
const TARGET_HEADER_NAME = "Target-Port"
const TARGET_PARAMETER_NAME = "bindToPort"
const PERFERS_HEADER_OVER_PARAM = true

// Configuration of a binding.
type Binding struct {
	// Net scheme used to transfert the request that will always override the default one.
	DefaultTransfertScheme *string `yaml:"transfert_scheme"`
	Enabled                bool    `yaml:"enabled"`
}
type Binder struct {
	TargetPort          int
	BindableConfig      *Binding
	OrginateFromRequest *http.Request
}

// Due to the binded port being possibly found on headers or parameters, this function is intended to choose between the value on headers
// or parameters if both are provided. Otherwise it searchs for a provided value in headers or parameters, if no value is provided the
// function will return an empty string.
func ExtractBindedPort(req *http.Request) (int, error) {
	var headerPort, headerErr = strconv.Atoi(req.Header.Get(TARGET_HEADER_NAME))
	var paramPort, paramErr = strconv.Atoi(req.URL.Query().Get(TARGET_PARAMETER_NAME))

	if headerErr != nil && paramErr != nil {
		return 0, errors.New("no viable binded parameter found")
	}

	if headerErr == nil && paramErr == nil {
		if PERFERS_HEADER_OVER_PARAM {
			return headerPort, nil
		} else {
			return paramPort, nil
		}
	} else {
		if headerErr == nil {
			return headerPort, nil
		} else {
			return paramPort, nil
		}
	}
}

// Creates a `Binder` from the request.
func GenerateBinder(req *http.Request, appropriateBinding *Binding) (*Binder, error) {
	// Extracts the port to use for this binder
	var port, err = ExtractBindedPort(req)

	// Determines if the binded configuration allows it to be used.
	if !appropriateBinding.Enabled {
		return nil, errors.New("binding" + strconv.Itoa(port) + "disabled")
	}

	if err != nil {
		return nil, errors.New("no " + TARGET_HEADER_NAME + " header found")
	}
	return &Binder{TargetPort: port, BindableConfig: appropriateBinding, OrginateFromRequest: req}, nil
}

// Binds a request to a specific port.
func (binder *Binder) BindToFromBinder() (*http.Response, error) {
	fmt.Printf("Binds regular from %s to ::%d", binder.OrginateFromRequest.RemoteAddr, binder.TargetPort)

	// Extracts `RequestURI` to build the new request URL
	binder.OrginateFromRequest.RequestURI = ""
	binder.OrginateFromRequest.URL.Host = fmt.Sprintf("localhost:%d", binder.TargetPort)

	// Determines which scheme to use due to the scheme override setting.
	if binder.BindableConfig.DefaultTransfertScheme != nil {
		binder.OrginateFromRequest.URL.Scheme = *binder.BindableConfig.DefaultTransfertScheme
	} else {
		binder.OrginateFromRequest.URL.Scheme = "http"
	}

	return http.DefaultClient.Do(binder.OrginateFromRequest)
}

// Create a SSE binding channel
func (binder *Binder) BindSseFromBinder(w http.ResponseWriter) error {
	fmt.Printf("Binds SSE from %s to ::%d", binder.OrginateFromRequest.RemoteAddr, binder.TargetPort)

	// Accept the incoming SSE request.
	incomingStream := MakeIncomingEventStream(w)

	if incomingStream.SupportEventStreams {
		// Make a request for the targeted server for an SSE connection
		outgoingStream, err := MakeOutgoingEventStreamFromRequest(binder.OrginateFromRequest, binder.BindableConfig)
		if err != nil {
			w.WriteHeader(500)
			return err
		}

		if err := outgoingStream.ReadEvents(
			func(readLine string) {
				incomingStream.WriteToBuffer(readLine)
				incomingStream.SendBuffer()
			}); err != nil {
			w.WriteHeader(500)
			return err
		}
	} else {
		// Sent if the client doesn't support event-streams.
		w.WriteHeader(400)
	}

	return nil
}
