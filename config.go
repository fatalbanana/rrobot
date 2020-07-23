package main

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Assertions []string          `hcl:"assertions"`
	Headers    map[string]string `hcl:"headers,optional"`
	Inputs     []string          `hcl:"inputs"`
	Name       string            `hcl:"name"`
}

func ReadConfig(configPath string) (error, Config) {
	var config Config
	val := hclsimple.DecodeFile(configPath, nil, &config)
	return val, config
}
