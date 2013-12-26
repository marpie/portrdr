package main

import (
	"encoding/json"
	"io/ioutil"
)

const (
	CONFIG_FILE = "portrdr.json"
)

type CertInfo struct {
	CertFile string `json:"certFile,omitempty"`
	KeyFile  string `json:"keyFile,omitempty"`
}

type RemoteAddrSsl struct {
	CertId     string `json:"certId,omitempty"`
	RemoteAddr string `json:"remoteAddr,omitempty"`
}

type Redirect struct {
	Info                 string              `json:"info,omitempty"`
	Protocol             string              `json:"protocol,omitempty"`
	LocalAddr            string              `json:"localAddr,omitempty"`
	SslSkipVerify        bool                `json:"sslSkipVerify,omitempty"`
	ApplicationProtocols []string            `json:"appProtos,omitempty"`
	RemoteAddr           string              `json:"remoteAddr,omitempty"`
	RemoteAddrs          []RemoteAddrSsl     `json:"remoteAddrs,omitempty"`
	Certs                map[string]CertInfo `json:"certs,omitempty"`
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
