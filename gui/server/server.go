package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"gocalcharger/api/service"
	"gocalcharger/gui/action"
	"gocalcharger/gui/tabs"
	"gocalcharger/server"
	"gocalcharger/server/check_permission"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"time"
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

// channels
var (
	CallbackChannel = make(chan action.ServerActionCallback, 1)
)

var serverGRPCServer *grpc.Server

func StartServer() {
	serverPort, _ = tabs.ServerPort.Get()
	serverPermitFiles, _ = tabs.ServerPermitFiles.Get()
	serverSSLEnabled, _ = tabs.ServerSSLEnabled.Get()
	serverSSLCert, _ = tabs.ServerSSLCert.Get()
	serverSSLKey, _ = tabs.ServerSSLKey.Get()
	serverCACert, _ = tabs.ServerCACert.Get()

	serveTarget := fmt.Sprintf(":%s", serverPort)
	listener, err := net.Listen("tcp", serveTarget)
	if err != nil {
		runCallbackError(serveTarget, action.ServerStartGRPCFailed, fmt.Sprintf("failed to listen: %v", err))
		return
	}

	err = check_permission.LoadPermission(serverPermitFiles)
	if err != nil {
		runCallbackError(serveTarget, action.ServerStartGRPCFailed, fmt.Sprintf("error loading permission:%v", err))
		return
	}
	// gRPC server.
	if serverSSLEnabled {
		// Mutual authentication.
		cert, err := tls.LoadX509KeyPair(serverSSLCert, serverSSLKey)
		if err != nil {
			runCallbackError(serveTarget, action.ServerStartGRPCFailed, fmt.Sprintf("can not load SSL credential:%v", err))
			return
		}
		certPool := x509.NewCertPool()
		credBytes, err := ioutil.ReadFile(serverCACert)
		if err != nil {
			runCallbackError(serveTarget, action.ServerStartGRPCFailed, fmt.Sprintf("can not load CA credential:%v", err))
			return
		}
		certPool.AppendCertsFromPEM(credBytes)
		cred := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
		})
		serverGRPCServer = grpc.NewServer(grpc.Creds(cred))
	} else {
		serverGRPCServer = grpc.NewServer()
	}
	service.RegisterGocalChargerServerServer(serverGRPCServer, &server.Server{})

	// reflection.Register(s)
	go func() {
		time.Sleep(time.Second)
		if serverGRPCServer == nil {
			return
		}
		log.Printf("gRPC server started[ServeTarget=%s]\n", serveTarget)
		CallbackChannel <- action.ServerActionCallback{
			CallbackName: action.ServerStartGRPCSuccess,
			CallbackArgs: action.ServerStartGRPCArgs{
				ServeTarget: serveTarget,
				Error:       nil,
			},
		}
	}()
	err = serverGRPCServer.Serve(listener)
	if err == nil {
		runCallbackError(serveTarget, action.ServerStartGRPCFailed, fmt.Sprintf("failed to serve: %v", err))
		return
	}
}

func runCallbackError(serveTarget string, errType action.ServerActionCallbackName, errString string) {
	log.Println(errString)
	CallbackChannel <- action.ServerActionCallback{
		CallbackName: errType,
		CallbackArgs: action.ServerStartGRPCArgs{
			ServeTarget: serveTarget,
			Error:       errors.New(errString),
		},
	}
}
