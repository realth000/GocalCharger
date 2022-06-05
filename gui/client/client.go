package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gocalcharger/client/file_download"
	"gocalcharger/client/say_hello"
	"gocalcharger/gui/action"
	"gocalcharger/gui/tabs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"log"
)

const (
	defaultName = "world gRPC!"
)

// Client configs
var (
	clientRemoteServerIP      string
	clientRemoteServerPort    string
	clientName                string
	clientSSLEnable           bool
	clientSSLCert             string
	clientSSLKey              string
	clientSSLCACert           string
	clientSSLMutualAuth       bool
	clientSSLDownloadFile     bool
	clientSSLDownloadFilePath string
)

// Channels
var (
	Channel         = make(chan action.ClientAction, 1)
	CallbackChannel = make(chan action.ClientActionCallback, 1)
)

func StartReceivingChannels() {
	go func() {
		for {
			select {
			case x := <-Channel:
				switch x.ActionName {
				case action.ClientSayHello:
					go SayHello(x.ActionArgs.(action.ClientSayHelloArgs).ClientName)
				case action.ClientDownloadFile:
					go DownloadFile(x.ActionArgs.(action.ClientDownloadFileArgs).FilePath)
				}
			}
		}
	}()
}

func SayHello(clientName string) {
	conn := initClient()
	if conn == nil {
		log.Fatalf("Nil connection in SayHello in client")
	}
	defer conn.Close()

	r, err := say_hello.SayHello(conn, updateClientName(clientName))
	if err != nil {
		log.Printf("error greeting: %v\n", err)
		CallbackChannel <- action.ClientActionCallback{
			CallbackName: action.ClientSayHelloFailed,
			CallbackArgs: action.ClientSayHelloCallbackArgs{
				ServerTarget: conn.Target(),
				Error:        err,
			},
		}
		return
	}
	log.Printf("successful greet: %s", r.Message)
	CallbackChannel <- action.ClientActionCallback{
		CallbackName: action.ClientSayHelloSuccess,
		CallbackArgs: action.ClientSayHelloCallbackArgs{
			ServerTarget: conn.Target(),
		},
	}
}

func DownloadFile(filePath string) {
	f := clientSSLDownloadFilePath
	if filePath != "" {
		f = filePath
	}
	conn := initClient()
	if conn == nil {
		log.Fatalf("Nil connection in Downloadfile in client")
	}
	defer conn.Close()

	file_download.DownloadFile(conn, updateClientName(""), f)
	log.Printf("download finish")
}

func updateClientName(name string) string {
	if name != "" {
		return name
	}
	if clientName != "" {
		return clientName
	}
	return defaultName
}

func initClient() *grpc.ClientConn {
	clientRemoteServerIP, _ = tabs.ClientRemoteServerIP.Get()
	clientRemoteServerPort, _ = tabs.ClientRemoteServerPort.Get()
	clientName, _ = tabs.ClientName.Get()
	clientSSLEnable, _ = tabs.ClientSSLEnable.Get()
	clientSSLCert, _ = tabs.ClientSSLCert.Get()
	clientSSLKey, _ = tabs.ClientSSLKey.Get()
	clientSSLCACert, _ = tabs.ClientSSLCACert.Get()
	clientSSLMutualAuth, _ = tabs.ClientSSLMutualAuth.Get()
	clientSSLDownloadFile, _ = tabs.ClientSSLDownloadFile.Get()
	clientSSLDownloadFilePath, _ = tabs.ClientSSLDownloadFilePath.Get()

	// Setup connection.
	if !clientSSLEnable {
		//conn, err := grpc.Dial(address, grpc.WithInsecure()) // deprecated
		c, err := grpc.Dial(fmt.Sprintf("%s:%s", clientRemoteServerIP, clientRemoteServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("can not dail %s:%s :%v", clientRemoteServerIP, clientRemoteServerPort, err)
		}
		return c
	} else {
		if clientSSLMutualAuth {
			// Mutual authentication.
			cert, err := tls.LoadX509KeyPair(clientSSLCert, clientSSLKey)
			if err != nil {
				log.Fatalf("can not load SSL credential:%v", err)
			}

			certPool := x509.NewCertPool()
			credBytes, err := ioutil.ReadFile(clientSSLCACert)
			if err != nil {
				log.Fatalf("can not load CA credential:%v", err)
			}

			certPool.AppendCertsFromPEM(credBytes)
			cred := credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
				//ServerName: "",
				RootCAs: certPool,
				// InsecureSkipVerify should be true to pass self-signed certificate.
				InsecureSkipVerify: true,
				// FIXME: Use some custom VerifyConnection to ensure handshake if using self-signed certificate.
				VerifyConnection: func(cs tls.ConnectionState) error {
					// TODO: Add verify test.
					opts := x509.VerifyOptions{
						DNSName:       cs.ServerName,
						Intermediates: x509.NewCertPool(),
					}
					for _, cert := range cs.PeerCertificates[1:] {
						opts.Intermediates.AddCert(cert)
					}
					_, err := cs.PeerCertificates[0].Verify(opts)
					return err
				},
			})
			c, err := grpc.Dial(fmt.Sprintf("%s:%s", clientRemoteServerIP, clientRemoteServerPort), grpc.WithTransportCredentials(cred))
			if err != nil {
				log.Fatalf("can not dail %s:%s :%v", clientRemoteServerIP, clientRemoteServerPort, err)
			}
			return c
		} else {
			cred, err := credentials.NewClientTLSFromFile(clientSSLCert, "")
			if err != nil {
				log.Fatalf("error ")
			}
			c, err := grpc.Dial(fmt.Sprintf("%s:%s", clientRemoteServerIP, clientRemoteServerPort), grpc.WithTransportCredentials(cred))
			if err != nil {
				log.Fatalf("can not dail %s:%s :%v", clientRemoteServerIP, clientRemoteServerPort, err)
			}
			return c
		}
	}
}
