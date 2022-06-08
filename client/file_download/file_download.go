package file_download

import (
	"context"
	"gocalcharger/api/service"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func DownloadFile(conn *grpc.ClientConn, name string, filePath string) {
	c := service.NewGocalChargerServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	request := service.DownloadFileRequest{
		ClientName: name,
		FileName:   filepath.Base(filePath),
		FilePath:   filePath,
	}
	r, err := c.DownloadFile(ctx, &request)
	if err != nil {
		log.Fatalf("error downloading file:%v", err)
	}

	// Delete old file
	_, err = os.Stat(request.FileName)
	if err == nil {
		os.Remove(request.FileName)
	}
	tmpFile, err := os.OpenFile(request.FileName, os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("can not save %s:%v", request.FileName, err)
	}

	var (
		fName    string
		size     int32
		progress int32
		total    int32
	)

	data, err := r.Recv()
	if err != nil {
		log.Fatalf("error receving file:%v", err)
	}
	fName = data.FileName
	size = data.FileSize
	progress++
	total = data.Total

	log.Printf("Download: %s[%3.0f%%][%d bytes]\033[1A", tmpFile.Name(), float32(data.Process)/float32(data.Total)*100, data.FileSize)
	_, err = tmpFile.Write(data.FilePart)
	if err != nil {
		log.Printf("error saving file %s:%v", request.FileName, err)
	}
	for {
		data, err := r.Recv()
		if err == io.EOF && progress == total {
			log.Printf("Download: %s[100%%][%d bytes]", fName, size)
			break
		}
		if err != nil {
			log.Fatalf("error receving file:%v", err)
		}
		progress++

		log.Printf("Download: %s[%3.0f%%][%d bytes]\033[1A", data.FileName, float32(data.Process)/float32(data.Total)*100, data.FileSize)
		_, err = tmpFile.Write(data.FilePart)
		if err != nil {
			log.Printf("error saving file %s:%v", request.FileName, err)
		}
	}
}
