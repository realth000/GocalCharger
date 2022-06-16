package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gocalcharger/api/service"
	"gocalcharger/server/check_permission"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"time"
)

type Server struct {
	service.UnimplementedGocalChargerServerServer
}

var (
	listener        net.Listener
	s               *grpc.Server
	ServerCloseChan = make(chan bool)
)

func (s *Server) SayHello(ctx context.Context, req *service.HelloRequest) (rsp *service.HelloReply, err error) {
	rsp = &service.HelloReply{Message: "Hello " + req.Name}
	log.Printf("Say Hello to %v\n", ctx)
	return rsp, nil
}

func (s *Server) DownloadFile(req *service.DownloadFileRequest, stream service.GocalChargerServer_DownloadFileServer) error {
	if !check_permission.CheckPathPermission(req.FilePath) {
		e := status.Error(codes.PermissionDenied, "denied to access this file")
		log.Printf("client=%s, %s", req.ClientName, e)
		return e
	}

	file, err := os.Open(req.FilePath)
	if err != nil {
		e := status.Error(codes.NotFound, "can not open this file"+err.Error())
		log.Printf("client=%s, %s", req.ClientName, e)
		return e
	}
	defer file.Close()
	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()
	const fileChunk = int32(1 * (1 << 20)) // 1 MB, change this to your requirement
	totalPartsNum := int32(math.Ceil(float64(fileSize) / float64(fileChunk)))
	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)
	for i := int32(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(float64(fileChunk), float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)
		file.Read(partBuffer)
		resp := &service.DownloadFileReply{
			FileName:  fileInfo.Name(),
			FileSize:  int32(fileInfo.Size()),
			FilePart:  partBuffer,
			Process:   i,
			Total:     totalPartsNum,
			FileChunk: fileChunk,
		}

		err = stream.SendMsg(resp)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}
	}
	return nil
}

func IsRunning() bool {
	time.Sleep(time.Millisecond * 500)
	if listener == nil {
		return false
	} else {
		return true
	}
}

func StartAndServe(port uint, permitConfig string) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	err = check_permission.LoadPermission(permitConfig)
	if err != nil {
		return err
	}
	// gRPC server.
	s := grpc.NewServer()
	service.RegisterGocalChargerServerServer(s, &Server{})
	// reflection.Register(s)
	fmt.Printf("gRPC serer running on %d\n", port)
	return s.Serve(listener)
}

func StartAndServeWithSSL(port uint, permitConfig string, cert string, key string, caCert string) error {
	var err error
	listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	err = check_permission.LoadPermission(permitConfig)
	if err != nil {
		return err
	}
	// gRPC server.
	// Mutual authentication.
	pair, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return err
	}
	certPool := x509.NewCertPool()
	credBytes, err := ioutil.ReadFile(caCert)
	if err != nil {
		log.Fatalf("can not load CA credential:%v", err)
	}
	certPool.AppendCertsFromPEM(credBytes)
	cred := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{pair},
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
	service.RegisterGocalChargerServerServer(s, &Server{})
	// reflection.Register(s)
	fmt.Printf("gRPC serer running on %d\n", port)
	go func() { err = s.Serve(listener) }()
	time.Sleep(time.Second)
	if !IsRunning() {
		return err
	} else {
		return nil
	}
}

func Stop() {
	if s != nil {
		s.Stop()
		s = nil
		ServerCloseChan <- true
	}
}
