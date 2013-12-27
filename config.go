package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

const (
	CONFIG_FILE = "portrdr.json"
)

var (
	ERR_HTTP_PORT_NOT_SET = errors.New("HTTP or HTTPS port is not set.")
)

type portRedirect struct {
	LocalAddr  string `json:"localAddr,omitempty"`
	RemoteAddr string `json:"remoteAddr,omitempty"`
}

type httpRedirect struct {
	Host      string `json:"host,omitempty"`
	Path      string `json:"path,omitempty"`
	RemoteURL string `json:"remoteURL,omitempty"`
	CertId    string `json:"certId,omitempty"`
}

type Config struct {
	HttpPort  uint16                     `json:"httpPort,omitempty"`
	HttpsPort uint16                     `json:"httpsPort,omitempty"`
	Tcp2Tcp   map[string]tcp2tcpRedirect `json:"tcp2tcp,omitempty"`
	Tcp2Udp   map[string]tcp2udpRedirect `json:"tcp2udp,omitempty"`
	Udp2Udp   map[string]tcp2udpRedirect `json:"udp2udp,omitempty"`
	Udp2Tcp   map[string]tcp2udpRedirect `json:"udp2tcp,omitempty"`
	Http      map[string]httpRedirect    `json:"http,omitempty"`
}

func LoadConfig(fileName string) (*Config, error) {
	// Read config file
	config_json, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, NewError("Reading config (%s): %v\n", fileName, err)
	}

	config := &Config{}

	// Parse config file
	if err = json.Unmarshal(config_json, config); err != nil {
		return nil, NewError("Parsing config (%s): %v\n", fileName, err)
	}

	// Check config
	if err = config.Check(); err != nil {
		return nil, err
	}

	return config, nil
}

func (cfg *Config) Check() error {
	if (len(cfg.Http) > 0) && (cfg.HttpPort+cfg.HttpsPort == 0) {
		return ERR_HTTP_PORT_NOT_SET
	}

	return nil
}

func (cfg *Config) Count() int {
	return len(cfg.Tcp2Tcp) + len(cfg.Tcp2Udp) + len(cfg.Udp2Tcp) + len(cfg.Udp2Udp) + len(cfg.Http)
}
