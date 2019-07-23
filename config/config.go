package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Role          string   `json:"role"`
	ServerAddr    string   `json:"server"`
	Tun           string   `json:"tun"`
	TunName       string   `json:"tunname"`
	PtcpInterface string   `json:"ptcpinterface"`
	Dns           string   `json:"dns"`
	Mtu           int      `json:"mtu"`
	Protocol      string   `json:"protocol"`
	Tokens        []string `json:"tokens"`
}

func NewConfig() *Config {
	return &Config{
		Role:          "Server",
		ServerAddr:    "127.0.0.1:12345",
		Tun:           "10.0.0.1/8",
		TunName:       "tun0",
		PtcpInterface: "eth0",
		Dns:           "8.8.8.8",
		Mtu:           1500,
		Protocol:      "tcp",
		Tokens:        []string{},
	}
}

func NewConfigFromFile(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cfg := NewConfig()
	err = cfg.Unmarshal(data)
	return cfg, err
}

func (cfg *Config) Unmarshal(data []byte) error {
	return json.Unmarshal(data, cfg)
}

func (cfg *Config) Marshal() ([]byte, error) {
	return json.Marshal(cfg)
}

func (cfg *Config) String() string {
	res, _ := cfg.Marshal()
	return string(res)
}
