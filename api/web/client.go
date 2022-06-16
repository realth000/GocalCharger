package web

import (
	"gocalcharger/client"
	"google.golang.org/grpc"
)

func ClientSayHello() (string, error) {
	var (
		c   *grpc.ClientConn
		err error
	)
	config := loadedConfigs.ClientConfig
	if config.SSL {
		c, err = client.DialSSL(config.ServerUrl, config.ServerPort, config.SSLCert)
		if err != nil {
			return "", err
		}
		if config.MutualAuth {
			c, err = client.DialSSLMutualAuth(config.ServerUrl, config.ServerPort, config.SSLCert, config.SSLKey, config.SSLCACert)
			if err != nil {
				return "", err
			}
		}
	} else {
		c, err = client.Dial(config.ServerUrl, config.ServerPort)
		if err != nil {
			return "", err
		}
	}
	return client.SayHello(c, config.ClientName)
}
