package server

import (
	"context"
	"fmt"
	"gocalcharger/api/service"
	"gocalcharger/server/check_permission"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math"
	"os"
)

type Server struct {
	service.UnimplementedGocalChargerServerServer
}

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
