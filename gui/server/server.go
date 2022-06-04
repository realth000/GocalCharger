package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gocalcharger/api/service"
	"gocalcharger/gui/tabs"
	"gocalcharger/server"
	"gocalcharger/server/check_permission"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
)

// Server configs
var (
	serverPort        string
	serverPermitFiles string
	serverSSLEnabled  bool
	serverSSLCert     string
	serverSSLKey      string
	serverCACert      string
)

func StartServer() {
	serverPort, _ = tabs.ServerPort.Get()
	serverPermitFiles, _ = tabs.ServerPermitFiles.Get()
	serverSSLEnabled, _ = tabs.ServerSSLEnabled.Get()
	serverSSLCert, _ = tabs.ServerSSLCert.Get()
	serverSSLKey, _ = tabs.ServerSSLKey.Get()
	serverCACert, _ = tabs.ServerCACert.Get()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	err = check_permission.LoadPermission(serverPermitFiles)
	if err != nil {
		log.Fatalf("error loading permission:%v\n", err)
	}
	// gRPC server.
	var s *grpc.Server
	if serverSSLEnabled {
		// Mutual authentication.
		cert, err := tls.LoadX509KeyPair(serverSSLCert, serverSSLKey)
		if err != nil {
			log.Fatalf("can not load SSL credential:%v", err)
		}
		certPool := x509.NewCertPool()
		credBytes, err := ioutil.ReadFile(serverCACert)
		if err != nil {
			log.Fatalf("can not load CA credential:%v", err)
		}
		certPool.AppendCertsFromPEM(credBytes)
		cred := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
		})
		s = grpc.NewServer(grpc.Creds(cred))
	} else {
		s = grpc.NewServer()
	}
	service.RegisterGocalChargerServerServer(s, &server.Server{})

	// reflection.Register(s)
	fmt.Printf("gRPC serer running on %s\n", serverPort)
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
