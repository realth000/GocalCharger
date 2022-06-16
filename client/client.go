package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"gocalcharger/client/say_hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
)

func Dial(ip string, port uint) (*grpc.ClientConn, error) {
	return grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func DialSSL(ip string, port uint, cert string) (*grpc.ClientConn, error) {
	cred, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		return nil, err
	}
	return grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(cred))
}

func DialSSLMutualAuth(ip string, port uint, cert string, key string, caCert string) (*grpc.ClientConn, error) {
	pair, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	credBytes, err := ioutil.ReadFile(caCert)
	if err != nil {
		return nil, err
	}

	certPool.AppendCertsFromPEM(credBytes)
	cred := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{pair},
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
	return grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(cred))
}

func SayHello(conn *grpc.ClientConn, name string) (string, error) {
	r, err := say_hello.SayHello(conn, name)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error greeting: %v", err))
	}
	return r.Message, nil
}

func DownloadFile(conn *grpc.ClientConn, name string) {

}