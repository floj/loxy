package config

import "github.com/hashicorp/hcl/v2/hclsimple"

type Config struct {
	Frontends []Frontend `hcl:"frontend,block"`
	Backends  []Backend  `hcl:"backend,block"`
}

func Load(path string) (*Config, error) {
	var config Config
	err := hclsimple.DecodeFile(path, nil, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
