package helpers

import (
	"johanmnto/epr/config"
	"johanmnto/epr/net"
	"net/http"
	"strconv"
)

// Determines if a request points to a binding in the configuration.
func PointsToKnownTarget(req *http.Request, config *config.EPRConfig) bool {
	var targetAsNumber, err = strconv.Atoi(req.Header.Get(net.TARGET_HEADER_NAME))
	if err != nil {
		return false
	}

	var _, ok = config.Bindings[targetAsNumber]
	return ok
}
