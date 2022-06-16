package main

import (
	"gocalcharger/client"
	"gocalcharger/client/config"
	"gocalcharger/client/file_download"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"sync"
)

const (
	defaultName = "world gRPC!"
)

// Command line flags.
var (
	flagConfigFile = kingpin.Flag("config-file", "Config file [*.toml].").String()

	flagServerUrl  = kingpin.Flag("server-ip", "Server IP.").Short('i').String()
	flagServerPort = kingpin.Flag("server-port", "Server port [0-65535].").Short('p').Uint()
	flagName       = kingpin.Flag("client-name", "Client name.").Short('n').String()

	flagSSL        = kingpin.Flag("ssl", "Use SSL in connecting with server. Use --no-ssl to disable ssl.").Default("true").Bool()
	flagSSLCert    = kingpin.Flag("cert", "SSL credential file[*.pem] path.").String()
	flagSSLKey     = kingpin.Flag("key", "SSL private key file[*.key] path.").String()
	flagSSLCACert  = kingpin.Flag("ca-cert", "SSL CA credential file[*.pem] path.").String()
	flagMutualAuth = kingpin.Flag("mutual-auth", "Mutual authentication in SSL handshake.").Default("true").Bool()

	cmdSayHello     = kingpin.Command("say-hello", "Say hello to server, used for testing.")
	cmdDownloadFile = kingpin.Command("download-file", "Download file from server.")
	flagFileName    = cmdDownloadFile.Flag("file-name", "Name of file to download.").String()
)

func loadConfigFile() {
	if *flagConfigFile == "" {
		return
	}
	err := config.LoadConfigFile(*flagConfigFile)
	cc := config.GetConfig()
	if err != nil {
		log.Fatalf("can not load config file:%v", err)
	}
	*flagServerUrl = cc.ServerUrl
	*flagServerPort = cc.ServerPort
	*flagName = cc.ClientName
	*flagSSL = cc.SSL
	*flagSSLCert = cc.SSLCert
	*flagSSLKey = cc.SSLKey
	*flagSSLCACert = cc.SSLCACert
	*flagMutualAuth = cc.MutualAuth
	*flagFileName = cc.DownloadFilePath
}

func checkFlag() {
	if *flagServerUrl == "" {
		*flagServerUrl = "localhost"
	}
	if *flagServerPort == 0 {
		log.Fatalln("Server port not set")
	} else if *flagServerPort > 65535 {
		log.Fatalf("Invalid port: %d\n", *flagServerPort)
	}
	switch kingpin.Parse() {
	case "say-hello":
	case "download-file":
		if *flagFileName == "" {
			log.Fatalf("Download path not set")
		}
	}

	if *flagSSL {
		if *flagSSLCert == "" {
			log.Fatalf("ssl enabled, but credential file[*.pem] not loaded")
		}
		if *flagMutualAuth {
			if *flagSSLKey == "" {
				log.Fatalf("ssl enabled, but private key file[*.key] not loaded")
			}
			if *flagSSLCACert == "" {
				log.Fatalf("mutual authentication enabled, but CA credential file[*.pem] not loaded")
			}
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

	// Setup connection.
	var conn *grpc.ClientConn
	if !*flagSSL {
		//conn, err := grpc.Dial(address, grpc.WithInsecure()) // deprecated
		c, err := client.Dial(*flagServerUrl, *flagServerPort)
		if err != nil {
			log.Fatalf("can not dail %s:%d :%v", *flagServerUrl, *flagServerPort, err)
		}
		conn = c
	} else {
		if *flagMutualAuth {
			c, err := client.DialSSLMutualAuth(*flagServerUrl, *flagServerPort, *flagSSLCert, *flagSSLKey, *flagSSLCACert)
			if err != nil {
				log.Fatalf("can not dail %s:%d :%v", *flagServerUrl, *flagServerPort, err)
			}
			conn = c
		} else {
			c, err := client.DialSSL(*flagServerUrl, *flagServerPort, *flagSSLCert)
			if err != nil {
				log.Fatalf("can not dail %s:%d :%v", *flagServerUrl, *flagServerPort, err)
			}
			conn = c
		}
	}
	defer conn.Close()
	// Contact the server and print out its response.
	var name string
	if *flagName == "" {
		name = defaultName
	} else {
		name = *flagName
	}
	switch kingpin.Parse() {
	case "say-hello":
		r, err := client.SayHello(conn, name)
		if err != nil {
			log.Fatalf("error greeting: %v\n", err)
		}
		log.Printf("successful greet: %s", r)
	case "download-file":
		wg := sync.WaitGroup{}
		wg.Add(1)
		go watchFileDownload(&wg)
		go file_download.DownloadFile(conn, name, *flagFileName)
		wg.Wait()
		finishDownloadWatchChan <- true
		log.Printf("download finish")
	}
}
