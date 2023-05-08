package config

import (
	"johanmnto/epr/net"
	"os"

	"gopkg.in/yaml.v2"
)

type EPRConfig struct {
	Server   net.Server           `yaml:"server"`
	Bindings map[int]*net.Binding `yaml:"bindings"`
}

const STATIC_CONFIG_PATH = "./config.epr.yaml"

func LoadAndParseConfiguration() EPRConfig {
	file, err := os.ReadFile(STATIC_CONFIG_PATH)

	if err != nil {
		panic(err)
	}

	var output EPRConfig
	err = yaml.Unmarshal(file, &output)
	if err != nil {
		panic(err)
	}

	return output
}
