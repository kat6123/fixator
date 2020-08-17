package config

import (
	"fixator/fixator"
	"fixator/handler"
	"flag"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type (
	Config struct {
		Fixator fixator.Config
		Service handler.Config
		Port    string `yaml:"port"`
		Host    string `yaml:"host"`
	}
)

func init() {
	flag.Parse()
}

var yamlPath = flag.String("yaml", "./config.yaml", "path to yaml file")

func Get() (*Config, error) {
	c := new(Config)

	content, err := ioutil.ReadFile(*yamlPath)
	if err != nil {
		return c, fmt.Errorf("config from yaml: read file %q: %v", *yamlPath, err)
	}

	err = yaml.Unmarshal(content, c)
	if err != nil {
		return c, fmt.Errorf("config from yaml: unmarshal: %v", err)
	}

	return c, nil
}
