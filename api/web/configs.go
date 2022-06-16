package web

import (
	clientConfig "gocalcharger/client/config"
	serverConfig "gocalcharger/server/config"
)

type Configs struct {
	ServerConfig serverConfig.ServerConfig
	ClientConfig clientConfig.ClientConfig
}

var loadedConfigs Configs

func getConfigs() Configs {
	loadedConfigs.ServerConfig = serverConfig.GetConfig()
	loadedConfigs.ClientConfig = clientConfig.GetConfig()
	return loadedConfigs
}

func reloadConfigs() Configs {
	// TODO: Update load config path.
	// TODO: Handle errors.
	_ = serverConfig.LoadConfigFile(`./tests/data/config/server.toml`)
	_ = clientConfig.LoadConfigFile(`./tests/data/config/client.toml`)
	return getConfigs()
}
