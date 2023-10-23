package net

import (
	"os"
	"gopkg.in/yaml.v2"
)

type NetConfig struct {
	Server       Server                 `yaml:"server"`
	Bindings     map[string]*Binding    `yaml:"bindings"`
}

func ParseConfigFrom(path string) NetConfig {
	file, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	var output NetConfig
	err = yaml.Unmarshal(file, &output)
	if err != nil {
		panic(err)
	}

	return output
}
