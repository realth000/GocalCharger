package file_download

import (
	"bytes"
	"context"
	"fmt"
	"gocalcharger/api/service"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
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
		log.Fatalf("error downloading file:%v\n", err)
	}

	b := new(bytes.Buffer)
	// Delete old file
	_, err = os.Stat(request.FileName)
	if err == nil {
		os.Remove(request.FileName)
	}
	var (
		fName    string
		size     int32
		progress int32
		total    int32
	)

	{
		data, err := r.Recv()
		if err != nil {
			fmt.Println(progress, total)
			log.Fatalf("error receving file:%v\n", err)

			return
		}
		fName = data.FileName
		size = data.FileSize
		progress++
		total = data.Total

		fmt.Printf("Download: %s[%3.0f%%][%d bytes]\n", data.FileName, float32(data.Process)/float32(data.Total)*100, data.FileSize)
		//progress_bar.UpdateProgress(request.FileName, int(100*size.Process/size.Total))
		b.Write(data.FilePart)
	}
	for {
		data, err := r.Recv()
		if err == io.EOF && progress == total {
			fmt.Printf("Download: %s[100%%][%d bytes]\n", fName, size)
			ioutil.WriteFile(request.FileName, b.Bytes(), 0755|os.ModeAppend)
			break
		}
		if err != nil {
			fmt.Println(progress, total)
			log.Fatalf("error receving file:%v\n", err)

			break
		}
		progress++

		fmt.Printf("Download: %s[%3.0f%%][%d bytes]\n", data.FileName, float32(data.Process)/float32(data.Total)*100, data.FileSize)
		//progress_bar.UpdateProgress(request.FileName, int(100*size.Process/size.Total))
		b.Write(data.FilePart)
	}
}
