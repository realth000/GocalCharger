package main

import (
	"gocalcharger/client/file_download"
	"log"
	"sync"
)

var (
	downloadChan            = &file_download.DownloadProgressChan
	finishDownloadWatchChan = make(chan bool)
)

func watchFileDownload(wg *sync.WaitGroup) {
	for {
		select {
		case x := <-*downloadChan:
			switch x.Status {
			case file_download.Downloading:
				log.Printf("Download: %s[%3.0f%%][%d bytes]\033[1A", x.Name, float32(x.Process)/float32(x.Total)*100, x.Size)
			case file_download.DownloadSuccess:
				log.Printf("Download: %s[100%%][%d bytes]", x.Name, x.Size)
				wg.Done()
			case file_download.DownloadFailed:
				log.Printf("fail to download %s:%s\n", x.Name, x.Err.Error())
				wg.Done()
			}
		case <-finishDownloadWatchChan:
			return
		}
	}
}
