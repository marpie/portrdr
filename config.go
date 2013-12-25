package main

import (
	"encoding/json"
	"io/ioutil"
)

const (
	CONFIG_FILE = "portrdr.json"
)

type Redirect struct {
	Info       string `json:"info,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	LocalAddr  string `json:"localAddr,omitempty"`
	RemoteAddr string `json:"remoteAddr,omitempty"`
}

type Config []Redirect

func LoadConfig(fileName string) (*Config, error) {
	// Read config file
	config_json, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, NewError("Reading config (%s): %v\n", fileName, err)
	}

	// Parse config file
	config := make(Config, 0)
	err = json.Unmarshal(config_json, &config)
	if err != nil {
		return nil, NewError("Parsing config (%s): %v\n", fileName, err)
	}

	return &config, nil
}
