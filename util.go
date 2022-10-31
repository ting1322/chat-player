package main

import (
	"strings"
	"path"
	"encoding/hex"
	"golang.org/x/crypto/blake2b"
)

// 先前的 python 版本行為，圖片存在 image 目錄，檔名 hash。
// 這邊使用相同的方式，取得跟以前一樣的檔名
func hashUrlFilename(image_url string) string {
	h, _ := blake2b.New(20, nil)
	h.Write([]byte(image_url))
	hashdata := h.Sum([]byte(image_url))
	if len(hashdata) > 20 {
		hashdata = hashdata[len(hashdata)-20:]
	}
	filename := "images/emoji-" + hex.EncodeToString(hashdata)
	if strings.HasSuffix(image_url, ".svg") {
		filename = filename + ".svg"
	} else {
		filename = filename + ".png"
	}
	return filename
}

// 檔名不要副檔名
func fileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(path.Ext(fileName))]
}