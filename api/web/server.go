package web

import "gocalcharger/server"

func startServer() error {
	config := loadedConfigs.ServerConfig
	if config.SSL {
		return server.StartAndServeWithSSL(config.Port, config.PermitFiles, config.SSLCert, config.SSLKey, config.SSLCACert)
	} else {
		return server.StartAndServe(config.Port, config.PermitFiles)
	}
}
