package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type ServerConfig struct {
	Port        uint   `toml:"port"`
	PermitFiles string `toml:"permit_files"`
	SSL         bool   `toml:"ssl"`
	SSLCert     string `toml:"ssl_cert"`
	SSLKey      string `toml:"ssl_key"`
	SSLCACert   string `toml:"ssl_ca_cert"`
}

var sc = ServerConfig{}

func GetConfig() ServerConfig {
	return sc
}

func LoadConfigFile(filePath string) error {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(configFile, &sc)
	return err
}
