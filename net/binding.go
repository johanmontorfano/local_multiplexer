package net

import (
	"fmt"
	"net/http"
	"strings"
)

// Header name used to determine which target is pointed.
const TARGET_HEADER_NAME = "Bind-To"

// Configuration of a server that can have requests binded to.
type Binding struct {
    Enabled             bool            `yaml:"enabled"`
    TargetPort          int             `yaml:"port"`

    // Automatically binds requests for specific paths. Thus, if `Binding` have
    // `auto_path` set to `/afs/`, every request that targets `/afs/` is 
    // requested to the `Binding`'s port.
    AutoBindPathRoot    *string         `yaml:"auto_path"`

    // Transfert protocol to use instead of the default one, this does not apply 
    // for SSE connections.
    TransfertProtocol   *string         `yaml:"protocol"`
}

// Returns the target `Binding` instance depending on the request and enabled
// bindings. This is useful to properly apply rules such as `AutoBindPathRoot`.
func GetAppropriateBinding(
    req *http.Request, 
    binds map[string]*Binding,
) *Binding {
    // Determines if one of the configured bindings have any auto rules defined
    // that can apply to this request.
    for bi := range binds {
        if binds[bi].AutoBindPathRoot != nil &&
            strings.HasPrefix(req.URL.Path, *binds[bi].AutoBindPathRoot) {
            return binds[bi]
        }
    }

    // Unwraps the port set in the `TARGET_HEADER_NAME` header. If no port is
    // set, then no appropriate binding cannot be found.
	var headerPort = req.Header.Get(TARGET_HEADER_NAME)

    return binds[headerPort]
}

// Binds a request to a specific port using HTTP or a custom transfert scheme.nf
func (bind *Binding) BindClassic(req *http.Request) (*http.Response, error) {
	fmt.Printf("Req. bind %s to ::%d\n", req.RemoteAddr, bind.TargetPort)

	req.RequestURI = ""
	req.URL.Host = fmt.Sprintf("localhost:%d", bind.TargetPort)
    req.URL.Scheme = "http"

    if bind.TransfertProtocol != nil {
        req.URL.Scheme = *bind.TransfertProtocol
    }

	return http.DefaultClient.Do(req)
}

// Create a SSE binding channel
func (bind *Binding) BindSSE(req *http.Request, w http.ResponseWriter) error {
	fmt.Printf("SSE. bind %s to ::%d\n", req.RemoteAddr, bind.TargetPort)

	// Accept the incoming SSE request.
	sse := MakeEventStream(w, req, *bind)

	if sse.SupportEventStreams {
		if err := sse.ReadEvents(func(readLine string) {
			sse.Add2Buffer(readLine)
			sse.SendBuffer()
		}); err != nil {
			w.WriteHeader(500)
			return err
		}
	} else {
		// Send this if the client doesn't support event-streams.
		w.WriteHeader(400)
	}

	return nil
}
