package helpers

import (
	"johanmnto/epr/config"
	"johanmnto/epr/net"
	"net/http"
)

// Determines if a request points to a binding in the configuration.
func PointsToKnownTarget(req *http.Request, config *config.EPRConfig) bool {
	var targetAsNumber, err = net.ExtractBindedPort(req)
	if err != nil {
		return false
	}

	var _, ok = config.Bindings[targetAsNumber]
	return ok
}
