package web

import (
	clientConfig "gocalcharger/client/config"
	serverConfig "gocalcharger/server/config"
)

type Configs struct {
	ServerConfig serverConfig.ServerConfig
	ClientConfig clientConfig.ClientConfig
}

func getConfigs() Configs {
	var c = Configs{}
	c.ServerConfig = serverConfig.GetConfig()
	c.ClientConfig = clientConfig.GetConfig()
	return c
}
