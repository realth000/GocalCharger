package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gocalcharger/api/service"
	"gocalcharger/api/web"
	"gocalcharger/server"
	"gocalcharger/server/check_permission"
	"gocalcharger/server/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"net"
)

var (
	flagConfigFile = kingpin.Flag("config-file", "Config file [*.toml].").String()

	flagPort        = kingpin.Flag("port", "gRPC running port.").Short('p').Uint()
	flagPermitFiles = kingpin.Flag("permit-files", "Load config file containing files permitted to access.").String()
	flagSSL         = kingpin.Flag("ssl", "Use SSL in connecting with server. Use --no-ssl to disable ssl.").Default("true").Bool()
	flagSSLCert     = kingpin.Flag("cert", "SSL credential file[*.pem] path.").String()
	flagSSLKey      = kingpin.Flag("key", "SSL private key file[*.key] path.").String()
	flagSSLCACert   = kingpin.Flag("ca-cert", "SSL CA credential file[*.pem] path.").String()
)

func loadConfigFile() {
	if *flagConfigFile == "" {
		return
	}
	err := config.LoadConfigFile(*flagConfigFile)
	sc := config.GetConfig()
	if err != nil {
		log.Fatalf("can not load config file:%v", err)
	}
	*flagPort = sc.Port
	*flagPermitFiles = sc.PermitFiles
	*flagSSL = sc.SSL
	*flagSSLCert = sc.SSLCert
	*flagSSLKey = sc.SSLKey
	*flagSSLCACert = sc.SSLCACert
}

func checkFlag() {
	if *flagPort == 0 {
		log.Fatalln("Port not set")
	} else if *flagPort > 65535 {
		log.Fatalf("Invalid port: %d\n", *flagPort)
	}
	if *flagSSL {
		if *flagSSLCert == "" {
			log.Fatalf("SSL enabled, but credential file[*.pem] not loaded")
		}
		if *flagSSLKey == "" {
			log.Fatalf("SSL enabled, but private key file[*.key] not loaded")
		}
		if *flagSSLCACert == "" {
			log.Fatalf("SSL enabled, but CA credential file[*.pem] not loaded")
		}
	}
}

func main() {
	// Setup flags.
	// Init flags? TODO: Do NOT use kingpin.Parse() in the line below.
	kingpin.Parse()
	// Load flags from config file.
	loadConfigFile()
	// Override flags with command line flags.
	kingpin.Parse()
	// Check if all flags legal.
	checkFlag()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *flagPort))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	err = check_permission.LoadPermission(*flagPermitFiles)
	if err != nil {
		log.Fatalf("error loading permission:%v\n", err)
	}
	// gRPC server.
	var s *grpc.Server
	if *flagSSL {
		// Mutual authentication.
		cert, err := tls.LoadX509KeyPair(*flagSSLCert, *flagSSLKey)
		if err != nil {
			log.Fatalf("can not load SSL credential:%v", err)
		}
		certPool := x509.NewCertPool()
		credBytes, err := ioutil.ReadFile(*flagSSLCACert)
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
		//if false {
		//	cred, err := credentials.NewServerTLSFromFile(argSSLCertPath, argSSLKeyPath)
		//	if err != nil {
		//		log.Fatalf("can not load SSL credential:%v", err)
		//	}
		//	s = grpc.NewServer(grpc.Creds(cred))
		//}

	} else {
		s = grpc.NewServer()
	}
	service.RegisterGocalChargerServerServer(s, &server.Server{})
	go web.RunServer()
	// reflection.Register(s)
	fmt.Printf("gRPC serer running on %d\n", *flagPort)
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
