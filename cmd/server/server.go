package main

import (
	"gocalcharger/server"
	"gocalcharger/server/config"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
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
	var err error
	if *flagSSL {
		err = server.StartAndServeWithSSL(*flagPort, *flagPermitFiles, *flagSSLCert, *flagSSLKey, *flagSSLCACert)
	} else {
		err = server.StartAndServe(*flagPort, *flagPermitFiles)
	}
	if err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
	<-server.ServerCloseChan
}
