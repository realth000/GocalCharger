package file_download

import (
	"bytes"
	"context"
	"gocalcharger/api/service"
	"gocalcharger/util/progress_bar"
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

	for {
		size, err := r.Recv()

		if err == io.EOF {
			log.Println("receive finish")
			ioutil.WriteFile(request.FileName, b.Bytes(), 0755|os.ModeAppend)
			break
		}
		if err != nil {
			log.Fatalf("error receving file:%v\n", err)
			break
		}
		progress_bar.UpdateProgress(request.FileName, int(100*size.Process/size.Total))
		b.Write(size.FilePart)
	}
}