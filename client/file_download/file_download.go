package file_download

import (
	"context"
	"gocalcharger/api/service"
	"google.golang.org/grpc"
	"io"
	"os"
	"path/filepath"
	"time"
)

type DownloadStatus int

const (
	Downloading DownloadStatus = iota
	DownloadSuccess
	DownloadFailed
)

type DownloadProgress struct {
	Name    string
	Status  DownloadStatus
	Process int32
	Total   int32
	Size    int32
	Err     error
}

var DownloadProgressChan = make(chan DownloadProgress, 1)

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
		//log.Fatalf("error downloading file:%v", err)
		DownloadProgressChan <- DownloadProgress{
			Name:    request.FileName,
			Status:  DownloadFailed,
			Process: 1,
			Total:   1,
			Size:    1,
			Err:     err,
		}
		return
	}

	// Delete old file
	_, err = os.Stat(request.FileName)
	if err == nil {
		os.Remove(request.FileName)
	}
	tmpFile, err := os.OpenFile(request.FileName, os.O_CREATE, 0644)
	if err != nil {
		//log.Fatalf("can not save %s:%v", request.FileName, err)
		DownloadProgressChan <- DownloadProgress{
			Name:    request.FileName,
			Status:  DownloadFailed,
			Process: 1,
			Total:   1,
			Size:    1,
			Err:     err,
		}
		return
	}

	var (
		fName    string
		size     int32
		progress int32
		total    int32
	)

	data, err := r.Recv()
	if err != nil {
		//log.Fatalf("error receving file:%v", err)
		DownloadProgressChan <- DownloadProgress{
			Name:    request.FileName,
			Status:  DownloadFailed,
			Process: 1,
			Total:   1,
			Size:    1,
			Err:     err,
		}
		return
	}
	fName = data.FileName
	size = data.FileSize
	progress++
	total = data.Total

	//log.Printf("Download: %s[%3.0f%%][%d bytes]\033[1A", tmpFile.Name(), float32(data.Process)/float32(data.Total)*100, data.FileSize)
	DownloadProgressChan <- DownloadProgress{
		Name:    fName,
		Status:  Downloading,
		Process: data.Process,
		Total:   data.Total,
		Size:    data.FileSize,
		Err:     nil,
	}
	_, err = tmpFile.Write(data.FilePart)
	if err != nil {
		//log.Printf("error saving file %s:%v", request.FileName, err)
		DownloadProgressChan <- DownloadProgress{
			Name:    fName,
			Status:  DownloadFailed,
			Process: data.Process,
			Total:   data.Total,
			Size:    data.FileSize,
			Err:     err,
		}
		return
	}
	for {
		data, err := r.Recv()
		if err == io.EOF && progress == total {
			//log.Printf("Download: %s[100%%][%d bytes]", fName, size)
			DownloadProgressChan <- DownloadProgress{
				Name:    fName,
				Status:  DownloadSuccess,
				Process: progress,
				Total:   total,
				Size:    size,
				Err:     err,
			}
			break
		}
		if err != nil {
			//log.Fatalf("error receving file:%v", err)
			DownloadProgressChan <- DownloadProgress{
				Name:    fName,
				Status:  DownloadFailed,
				Process: progress,
				Total:   total,
				Size:    size,
				Err:     err,
			}
			return
		}
		progress++

		//log.Printf("Download: %s[%3.0f%%][%d bytes]\033[1A", data.FileName, float32(data.Process)/float32(data.Total)*100, data.FileSize)
		DownloadProgressChan <- DownloadProgress{
			Name:    fName,
			Status:  Downloading,
			Process: data.Process,
			Total:   data.Total,
			Size:    data.FileSize,
			Err:     err,
		}
		_, err = tmpFile.Write(data.FilePart)
		if err != nil {
			//log.Printf("error saving file %s:%v", request.FileName, err)
			DownloadProgressChan <- DownloadProgress{
				Name:    fName,
				Status:  DownloadFailed,
				Process: progress,
				Total:   total,
				Size:    size,
				Err:     err,
			}
			return
		}
	}
}
