package net

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// Header name used to determine which target is pointed.
const TARGET_HEADER_NAME = "Target-Port"

// Configuration of a binding.
type Binding struct {
	MaxResWaitTime int `yaml:"max_res_wait_time"`
	// Net scheme used to transfert the request that will always override the default one.
	DefaultTransfertScheme *string `yaml:"transfert_scheme"`
}
type Binder struct {
	TargetPort          int
	BindableConfig      *Binding
	OrginateFromRequest *http.Request
}

// Creates a `Binder` from the request.
func GenerateBinder(req *http.Request, appropriateBinding *Binding) (*Binder, error) {
	var port, err = strconv.Atoi(req.Header.Get(TARGET_HEADER_NAME))
	if err != nil {
		return nil, errors.New("no " + TARGET_HEADER_NAME + " header found")
	}
	return &Binder{TargetPort: port, BindableConfig: appropriateBinding, OrginateFromRequest: req}, nil
}

// Binds a request to a specific port.
func (binder *Binder) BindToFromBinder() (*http.Response, error) {
	// Extracts `RequestURI` to build the new request URL
	binder.OrginateFromRequest.RequestURI = ""
	binder.OrginateFromRequest.URL.Host = fmt.Sprintf("localhost:%d", binder.TargetPort)

	// Determines which scheme to use due to the scheme override setting.
	if binder.BindableConfig.DefaultTransfertScheme != nil {
		binder.OrginateFromRequest.URL.Scheme = *binder.BindableConfig.DefaultTransfertScheme
	} else {
		binder.OrginateFromRequest.URL.Scheme = "http"
	}

	// println(binder.OrginateFromRequest.URL.RequestURI())
	return http.DefaultClient.Do(binder.OrginateFromRequest)
}
