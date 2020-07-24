package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type RealConfig struct {
	Tests []Config `hcl:"test,block"`
}

type Config struct {
	Assertions []string          `hcl:"assertions"`
	Headers    map[string]string `hcl:"headers,optional"`
	Inputs     []string          `hcl:"inputs"`
	Name       string            `hcl:"name,label"`
}

func ReadConfig(configPath string) (RealConfig, error) {
	var config RealConfig
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(configPath)
	if diags.HasErrors() {
		return config, fmt.Errorf("parse error: %s", diags.Error())
	}
	diags = gohcl.DecodeBody(file.Body, nil, &config)
	if diags.HasErrors() {
		return config, fmt.Errorf("decode error: %s", diags.Error())
	}
	return config, nil
}
