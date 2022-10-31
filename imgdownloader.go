package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ImgDownloader interface {
	Download(localPath, url string) string
}

type HttpImgDownloader struct {}

func (me *HttpImgDownloader) Download(localPath, url string) string {
	_, err := os.Stat(localPath)
	if err == nil {
		return url
	}
	if strings.HasPrefix(url, "//") {
		url = "https:" + url
	}
	_, err = os.Stat(filepath.Dir(localPath))
	if os.IsNotExist(err) {
		os.Mkdir(filepath.Dir(localPath), os.ModePerm)
	}
	
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Opaque = req.URL.Path
			return nil
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Println("靠北，下載表情符號失敗了: " + url)
		return url
	}
	defer resp.Body.Close()
	file, err := os.Create(localPath)
	if err != nil {
		return url
	}
	defer file.Close()
	io.Copy(file, resp.Body)
	log.Println("download emilji: " + url)
	log.Println("save to: " + localPath)
	return localPath
}