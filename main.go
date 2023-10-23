package main

import (
	"flag"
	"fmt"
	"io"
	jnet "johanmnto/multiplexer/net"
	"net/http"
	"strings"
	"sync"
)

func main() {
    config_path := flag.String(
        "config", 
        "./multiplexer.yaml", 
        "Path to the configuration file that has to be loaded")

    flag.Parse()

    config := jnet.ParseConfigFrom(*config_path)
    wait_group := new(sync.WaitGroup)

	server := jnet.MakeServer(&config)
	server.MakeHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")

        // When client's requires an OPTIONS requst, automatically answer 
        // positively.
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
		    return
        }

        // Determines if a binding is configured for the incoming request.
        binding := jnet.GetAppropriateBinding(r, config.Bindings)

		if (binding != nil && binding.Enabled) {
			// Here, the handler will automatically determines if the request 
            // should be handled as a standard http requestor as an event stream.
			// Note that the handler will only forward to SSE handlers request 
            // which has the `Accept` header set as `text/event-stream`.
			if r.Header.Get("Accept") == "text/event-stream" {
                binding.BindSSE(r, w)
			} else {
				// Transfers the request to the binded port
				response, err := binding.BindClassic(r)

				if err != nil {
					println(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					// If the transfert had succeded and a response has been 
                    // given, the response is copied to be sent back to the 
                    // client.
					body, err := io.ReadAll(response.Body)

					if err != nil {
						println(err.Error())
						w.WriteHeader(http.StatusInternalServerError)
					} else {
						w.WriteHeader(response.StatusCode)
						for headerName, headerValue := range response.Header {
							w.Header().Add(
                                headerName, 
                                strings.Join(headerValue, ", "),
                            )
						}
						w.Write(body)
					}
				}
			}
		} else { 
			w.WriteHeader(http.StatusNotFound)
            w.Write([]byte("NOT FOUND"))
		}
	})

	// Starts both unsecure and secure server, the secure server is started only 
    // if an HTTPS port is provided.
    //
    // This schema for starting servers allows for both processes to not provoke
    // a global app fail and to run concurrently. 
    // Thus, if one of the two servers is taken down, it can be restarted.
	wait_group.Add(1)
	go func() {
        for {
		    if err := server.ServeUnsecure(); err != nil {
			    panic(fmt.Sprintf("HTTP server failure: %s", err.Error()))
		    }
        }
	}()
	if config.Server.HttpsPort != nil {
		wait_group.Add(1)
		go func() {
            for {
			    if err := server.ServeSecure(); err != nil {
				    panic(fmt.Sprintf("HTTPS server failure: %s", err.Error()))
			    }
            }
		}()
	}

    wait_group.Wait()
}
