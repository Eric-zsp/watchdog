package files

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// const (
// 	chunkSize     = 1024 * 1024 // 1 MB
// 	progressWidth = 50
// 	// bufferSize    = 1024 * 1024 // 更新进度的缓冲字节数
// )

// func DownloadFile(url string, destPath string) error {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return fmt.Errorf("failed to start download: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("download failed with status: %s", resp.Status)
// 	}

// 	out, err := os.Create(destPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create output file: %v", err)
// 	}
// 	defer out.Close()

// 	totalSize := resp.ContentLength
// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	go func() {
// 		defer wg.Done()
// 		progressBar(totalSize, resp.Body)
// 	}()

// 	_, err = io.Copy(out, resp.Body)
// 	if err != nil && !strings.Contains(err.Error(), "unexpected EOF") {
// 		return fmt.Errorf("error writing to output file: %v", err)
// 	}

// 	wg.Wait()
// 	fmt.Println("\nDownload completed")
// 	return nil
// }

// func progressBar(total int64, body io.Reader) {
// 	reader := &progressReader{rd: body, total: total}
// 	buffer := make([]byte, chunkSize)
// 	for {
// 		n, err := reader.Read(buffer)
// 		if err != nil && err != io.EOF {
// 			gologs.GetLogger("default").Sugar().Errorf("Error reading data: %v", err)
// 			break
// 		}
// 		if n > 0 {
// 			if reader.readed >= reader.total*10/100 && reader.readed < reader.total*20/100 {
// 				updateProgress(reader.readed, reader.total)
// 			} else if reader.readed >= reader.total*20/100 && reader.readed < reader.total*30/100 {
// 				updateProgress(reader.readed, reader.total)
// 			} else if reader.readed >= reader.total*30/100 && reader.readed < reader.total*40/100 {
// 				updateProgress(reader.readed, reader.total)
// 			} else if reader.readed >= reader.total*40/100 && reader.readed < reader.total*50/100 {
// 				updateProgress(reader.readed, reader.total)
// 			} else if reader.readed >= reader.total*50/100 && reader.readed < reader.total*60/100 {
// 				updateProgress(reader.readed, reader.total)
// 			} else if reader.readed >= reader.total*60/100 && reader.readed < reader.total*70/100 {
// 				updateProgress(reader.readed, reader.total)
// 			} else if reader.readed >= reader.total*70/100 && reader.readed < reader.total*80/100 {
// 				updateProgress(reader.readed, reader.total)
// 			} else if reader.readed >= reader.total*80/100 && reader.readed < reader.total*90/100 {
// 				updateProgress(reader.readed, reader.total)
// 			} else if reader.readed >= reader.total*90/100 {
// 				updateProgress(reader.readed, reader.total)
// 			}
// 		} else {
// 			break
// 		}
// 	}
// }

// type progressReader struct {
// 	rd     io.Reader
// 	total  int64
// 	readed int64
// }

// func (pr *progressReader) Read(p []byte) (n int, err error) {
// 	n, err = pr.rd.Read(p)
// 	pr.readed += int64(n)
// 	return n, err
// }

// func updateProgress(current, total int64) {
// 	progress := int(float64(current) / float64(total) * float64(progressWidth))
// 	bar := strings.Repeat("=", progress) + strings.Repeat(" ", progressWidth-progress)
// 	gologs.GetLogger("default").Sugar().Infof("\rProgress: [%s] %.2f%%", bar, float64(current)/float64(total)*100)
// 	// 添加这一行，确保进度条在100%之前持续更新
// 	if current >= total {
// 		fmt.Printf("\n")
// 	}
// }

const (
	concurrency = 10          // 并发线程数
	chunkSize   = 1024 * 1024 // 每个线程下载的大小，这里是1MB
)

type downloadProgress struct {
	total       int64
	downloaded  int64
	mu          sync.Mutex
	progressBar *sync.Cond
}

func (dp *downloadProgress) updateDownloaded(newDownloaded int64) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.downloaded += newDownloaded
	dp.progressBar.Broadcast()
}

func (dp *downloadProgress) printProgress() {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	progress := float64(dp.downloaded) * 100 / float64(dp.total)
	bar := strings.Repeat("=", int(dp.downloaded/dp.total*50)) + strings.Repeat(" ", 50-int(dp.downloaded/dp.total*50))
	fmt.Printf("\rProgress: [%s] %.2f%%", bar, progress)
}

func DownloadFile(url string, destPath string) error {
	resp, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("failed to get file size: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	contentLength, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if contentLength <= 0 {
		return fmt.Errorf("invalid content length")
	}

	dp := &downloadProgress{
		total:       contentLength,
		progressBar: sync.NewCond(&sync.Mutex{}),
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	wg.Add(concurrency)

	blockSize := contentLength / int64(concurrency)
	offsets := make([]int64, concurrency)
	for i := 0; i < concurrency; i++ {
		offsets[i] = int64(i) * blockSize
		go downloadChunk(&wg, url, offsets[i], file, blockSize, contentLength, dp)
	}

	go dp.printProgressLoop()

	wg.Wait()
	fmt.Println("\nDownload completed")
	return nil
}

func (dp *downloadProgress) printProgressLoop() {
	for {
		dp.progressBar.L.Lock()
		for dp.downloaded < dp.total {
			dp.progressBar.Wait()
		}
		dp.progressBar.L.Unlock()
		dp.printProgress()
	}
}

func downloadChunk(wg *sync.WaitGroup, url string, offset int64, file *os.File, blockSize int64, total int64, dp *downloadProgress) {
	defer wg.Done()

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return
	}
	request.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, offset+blockSize-1))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Error downloading chunk: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		log.Printf("Unexpected status code: %v", resp.StatusCode)
		return
	}

	buffer := make([]byte, blockSize)
	for {
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			log.Printf("Error reading chunk: %v", err)
			return
		}
		if n > 0 {
			_, err = file.WriteAt(buffer[:n], offset)
			if err != nil {
				log.Printf("Error writing to file: %v", err)
				return
			}
			dp.updateDownloaded(int64(n))
		}
		if err == io.EOF {
			break
		}
	}
}
