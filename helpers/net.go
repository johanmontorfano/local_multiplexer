package helpers

import (
	"johanmnto/epr/config"
	"johanmnto/epr/net"
)

// Builds a `Server` from `EPRConfig`
func MakeServer(config *config.EPRConfig) net.Server {
	return net.Server{
		HttpPort:      config.Server.HttpPort,
		HttpsPort:     config.Server.HttpsPort,
		HttpsKeyPath:  config.Server.HttpsKeyPath,
		HttpsCertPath: config.Server.HttpsCertPath,
	}
}
