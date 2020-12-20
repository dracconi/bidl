package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"
	"os"

	"github.com/dracconi/bidl/downloaders"
)

var finished bool

func main() {
	start := time.Now()
	downloaders.InitRules()
	var wg sync.WaitGroup
	finished = false
	var bulk_urls []string
	var prefix string
	for _, v := range os.Args[1:] {
		if (v[:2] == "-o") {
			prefix = v[2:]
		} else {
			bulk_urls = append(bulk_urls, v)			
		}
	}
	fmt.Println(prefix)
	files := make(chan downloaders.RemoteImage, 256)
	go workerPool(files, &wg, prefix)
	count := expandUrls(bulk_urls, files)
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Finished %d in %fs. Averaging %dms.", count, elapsed.Seconds(), elapsed.Milliseconds()/count)
}

func expandUrls(urls []string, files chan downloaders.RemoteImage) int64 {
	var count int64 = 0
	for _, v := range urls {
		vurl, _ := url.Parse(v)
		vurls, _ := downloaders.GetImUrls(vurl)
		for _, w := range vurls {
			fmt.Println(w)
			count = count + 1
			files <- w
		}
	}
	finished = true
	return count
}
