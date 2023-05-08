package net

import (
	"fmt"
	"net/http"
)

type Server struct {
	HttpPort      int     `yaml:"http_port"`
	HttpsPort     *int    `yaml:"https_port"`
	HttpsKeyPath  *string `yaml:"https_key_path"`
	HttpsCertPath *string `yaml:"https_cert_path"`

	mux *http.ServeMux
}

// Allocates a new server mux to the server that will be used when starting the server.
func (srv *Server) MakeHandler(handler func(http.ResponseWriter, *http.Request)) {
	srv.mux = http.NewServeMux()
	srv.mux.HandleFunc("/", handler)
}

func (srv *Server) ServeUnsecure() error {
	fmt.Println("Listening for HTTP requests on port", srv.HttpPort)
	return http.ListenAndServe(fmt.Sprintf(":%d", srv.HttpPort), srv.mux)
}

func (srv *Server) ServeSecure() error {
	fmt.Println("Listening for HTTPS requests on port", *srv.HttpsPort)
	return http.ListenAndServeTLS(fmt.Sprintf(":%d", *srv.HttpsPort), *srv.HttpsCertPath, *srv.HttpsKeyPath, srv.mux)
}
