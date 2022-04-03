package parser

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

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

type Condition struct {
	Type  string `yaml:"type"`
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
	State string `yaml:"state"`
}

type Delay struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type Response struct {
	StatusCode int `yaml:"statusCode"`
	Headers    []struct {
		Key   string `yaml:"key"`
		Value string `yaml:"value"`
	} `yaml:"headers"`
	Body      string    `yaml:"body"`
	Delay     Delay     `yaml:"delay"`
	Condition Condition `yaml:"condition"`
}

type Route struct {
	Method   string     `yaml:"method"`
	Endpoint string     `yaml:"endpoint"`
	Response []Response `yaml:"responses"`
}

func Parse_YAML(filename string) (Config, error) {
	var config Config

	var data []byte
	var err error
	var response *http.Response

	if strings.HasPrefix(filename, "http") {
		response, err = http.Get(filename)
		if err != nil {
			return config, err
		}
		defer response.Body.Close()
		n, e2 := io.ReadAll(response.Body)
		if e2 != nil {
			return config, e2
		}
		data = n
	} else {
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			return config, nil
		}
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
