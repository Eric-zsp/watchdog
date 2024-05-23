package files

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	gologs "github.com/cn-joyconn/gologs"
	"github.com/cn-joyconn/goutils/filetool"
)

// @title   保存文件
//
//	path 文件路径
//	bs 保存内容(byte)
func SaveFileByes(fileName string, bs []byte) (string, error) {

	selfDir := filetool.SelfDir()
	fullPath := path.Join(selfDir, fileName)
	_, err := filetool.WriteBytesToFile(fullPath, bs)
	if err != nil {
		return "", err
	}
	return fullPath, nil

}

// func CopyFile(dstName, srcName string) (written int64, err error) {
// 	src, err := os.Open(srcName)
// 	if err != nil {
// 		return
// 	}
// 	defer src.Close()
// 	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
// 	if err != nil {
// 		return
// 	}
// 	defer dst.Close()
// 	return io.Copy(dst, src)
// }

// // copyDir 复制整个目录
// func CopyDir(src, dst string) error {
// 	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		relPath, err := filepath.Rel(src, path)
// 		if err != nil {
// 			return err
// 		}

// 		dstPath := filepath.Join(dst, relPath)

// 		if info.IsDir() {
// 			return os.MkdirAll(dstPath, info.Mode())
// 		} else {
// 			// 检查是否可以写入，如果不能，尝试获取写权限
// 			// if err := os.Chmod(dstPath, info.Mode()); err != nil && !os.IsExist(err) {
// 			// 	return err
// 			// }
// 			_, err = CopyFile(path, dstPath)
// 			return err
// 		}
// 	})
// 	return err
// }

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// 如果路径不存在或者有其他错误，这里会返回相应的错误信息
		gologs.GetLogger("default").Sugar().Error("Error accessing path: %v\n", err)
		return false
	}
	return info.IsDir()
}
func DownLoadFile2(url string, destPath string) error {
	// 创建一个HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to start download: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// 打开目标文件，准备写入
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer out.Close()

	// 将响应体的内容复制到输出文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing to output file: %v", err)
	}

	fmt.Printf("Successfully downloaded %s\n", destPath)
	return nil
}
func DownLoadFile1(durl string, savePath string) {
	_, err := url.ParseRequestURI(durl)
	if err != nil {
		panic("网址错误")
	}

	// filename := path.Base(uri.Path)
	// gologs.GetLogger("default").Info("[*] Filename " + filename)

	client := http.DefaultClient
	client.Timeout = time.Second * 60 //设置超时时间
	resp, err := client.Get(durl)
	if err != nil {
		panic(err)
	}
	if resp.ContentLength <= 0 {
		gologs.GetLogger("default").Error("[*] Destination server does not support breakpoint download.")
	}
	raw := resp.Body
	defer raw.Close()
	reader := bufio.NewReaderSize(raw, 1024*32)

	file, err := os.Create(savePath)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)

	buff := make([]byte, 32*1024)
	written := 0
	go func() {
		for {
			nr, er := reader.Read(buff)
			if nr > 0 {
				nw, ew := writer.Write(buff[0:nr])
				if nw > 0 {
					written += nw
				}
				if ew != nil {
					err = ew
					break
				}
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
		}
		if err != nil {
			panic(err)
		}
	}()

	spaceTime := time.Second * 1
	ticker := time.NewTicker(spaceTime)
	lastWtn := 0
	stop := false

	for {
		select {
		case <-ticker.C:
			// speed := written - lastWtn
			// gologs.GetLogger("default").Sugar().Info("[*] Speed %s / %s \n", bytesToSize(speed), spaceTime.String())
			if written-lastWtn <= 0 {
				ticker.Stop()
				stop = true
				break
			}
			lastWtn = written
		}
		if stop {
			break
		}
	}
}
func bytesToSize(length int) string {
	var k = 1024 // or 1024
	var sizes = []string{"Bytes", "KB", "MB", "GB", "TB"}
	if length == 0 {
		return "0 Bytes"
	}
	i := math.Floor(math.Log(float64(length)) / math.Log(float64(k)))
	r := float64(length) / math.Pow(float64(k), i)
	return strconv.FormatFloat(r, 'f', 3, 64) + " " + sizes[int(i)]
}

// CopyFile copies the contents of the file named src to the file named dst.
// If dst doesn't exist, it's created. If src and dst are the same file, CopyFile returns nil.
func CopyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	return err
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must not.
func CopyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
