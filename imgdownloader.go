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
	if filename, exist := IsLocalFileExist(localPath); exist {
		return filename
	}
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
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("靠北，下載表情符號失敗了: " + url)
		return url
	}
	if filepath.Ext(localPath) != ".svg" {
		ext := GetExtension(http.DetectContentType(data))
		localPath = fileNameWithoutExt(localPath) + ext
	}
	file, err := os.Create(localPath)
	if err != nil {
		return url
	}
	defer file.Close()
	file.Write(data)
	log.Println("download emilji: " + url)
	log.Println("save to: " + localPath)
	return localPath
}

func GetExtension(contentType string) string {
	switch contentType {
	case "image/webp": return ".webp"
	case "image/jpeg": return ".jpg"
	case "image/gif": return ".gif"
	case "image/bmp": return ".bmp"
	case "image/png": return ".png"
	default: return ".png"
	}
}

func IsLocalFileExist(localPath string) (filename string, exist bool) {
	_, err := os.Stat(localPath)
	if err == nil {
		return localPath, true
	}
	p := fileNameWithoutExt(localPath) + ".webp"
	if _, err := os.Stat(p); err == nil {
		return p, true
	}
	p = fileNameWithoutExt(localPath) + ".svg"
	if _, err := os.Stat(p); err == nil {
		return p, true
	}
	p = fileNameWithoutExt(localPath) + ".jpg"
	if _, err := os.Stat(p); err == nil {
		return p, true
	}
	p = fileNameWithoutExt(localPath) + ".gif"
	if _, err := os.Stat(p); err == nil {
		return p, true
	}
	p = fileNameWithoutExt(localPath) + ".bmp"
	if _, err := os.Stat(p); err == nil {
		return p, true
	}
	p = fileNameWithoutExt(localPath) + ".png"
	if _, err := os.Stat(p); err == nil {
		return p, true
	}
	return localPath, false
}