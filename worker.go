package main

import (
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/dracconi/bidl/downloaders"
)

func doesFileNotExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsExist(err)
}

func downloadFile(dest *os.File, src string) (int64, error) {
	resp, err := http.Get(src)
	if err != nil {
		return 0, err
	}
	return dest.ReadFrom(resp.Body)
}

func worker(ch chan int, file downloaders.RemoteImage, wg *sync.WaitGroup, prefix string) {
	if doesFileNotExist(prefix+file.Local) {
		err := os.MkdirAll(path.Dir(prefix + file.Local), 0777)
		if err != nil {
			panic(err)
		}
		f, _ := os.Create(prefix+file.Local)
		downloadFile(f, file.Remote)
		f.Close()
	}
	wg.Done()
	ch <- 1
}

func workerLoop(ch chan int, files chan downloaders.RemoteImage, wg *sync.WaitGroup, prefix string) {
	br := true
	for br {
		select {
		case f := <-files:
			<-ch
			wg.Add(1)
			go worker(ch, f, wg, prefix)
		default:
			if finished {
				br = false
			}
		}
	}
}

func workerPool(files chan downloaders.RemoteImage, wg *sync.WaitGroup, prefix string) {
	ch := make(chan int, 129)
	for i := 0; i < 128; i++ {
		ch <- 1
	}
	workerLoop(ch, files, wg, prefix)
}
