package parser

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name     string  `yaml:"name"`
	Hostname string  `yaml:"hostname"`
	Port     string  `yaml:"port"`
	Options  Options `yaml:"options"`
	Route    []Route `yaml:"routes"`
}

type Options struct {
	AccessControlAllowOrigin      string `yaml:"accessControlAllowOrigin"`
	AccessControlAllowCredentials string `yaml:"accessControlAllowCredentials"`
	AccessControlAllowHeaders     string `yaml:"accessControlAllowHeaders"`
	AccessControlAllowMethods     string `yaml:"accessControlAllowMethods"`
}

type Route struct {
	Method   string `yaml:"method"`
	Endpoint string `yaml:"endpoint"`
	Response struct {
		StatusCode int `yaml:"statusCode"`
		Headers    []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		} `yaml:"headers"`
		Body string `yaml:"body"`
	} `yaml:"response"`
}

func Parse_YAML(filename string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
