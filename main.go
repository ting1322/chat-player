package main

import (
	_ "embed"
	"errors"
	"log"
	"os"
	"strings"

	"path/filepath"

	"bufio"
)

//go:embed template.htm.in
var templateHtm string

//go:embed play-live-chat.js
var playlivechatjs string

func main() {
	parseCommandline()

	fileInfo, err := os.Stat(optionPath)
	if err != nil {
		log.Fatalln(err)
	}
	if fileInfo.IsDir() {
		if len(optionChatJson) != 0 ||
			len(optionSetList) != 0 {
			log.Fatalln("not support option in directory mode")
		}
		webms, err := filepath.Glob(optionPath + "/*.webm")
		if err != nil {
			log.Fatalln(err)
		}
		mp4s, err := filepath.Glob(optionPath + "/*.mp4")
		if err != nil {
			log.Fatalln(err)
		}
		videoList := append(webms, mp4s...)
		if len(videoList) == 0 {
			log.Fatalln("not found ant video file")
		}
		for _, video := range videoList {
			err = processVideo(video)
			if err != nil &&
				err.Error() != not_found_json {
				log.Fatalln(err)
			}
		}
		
	} else {
		err = processVideo(optionPath)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

const not_found_json string = "not found .live_chat.json file"

func processVideo(filePath string) error {
	var chatJsonFilename string = optionChatJson
	if len(chatJsonFilename) == 0 {
		chatJsonFilename = fileNameWithoutExt(filePath) + ".live_chat.json"
	}
	_, err := os.Stat(chatJsonFilename)
	if err != nil {
		return errors.New(not_found_json)
	}
	_, err = os.Stat(filePath)
	if err != nil {
		return err
	}
	var outputFilename string = optionOutputName
	var outDir string = optionOutDir
	if len(optionOutputName) == 0 {
		outputFilename = filePath + ".htm"
	}
	if len(outDir) > 0 {
		outputFilename = filepath.Join(outDir, filepath.Base(outputFilename))
	} else {
		outDir = filepath.Dir(outputFilename)
	}
	videoPathInHtm, err := filepath.Rel(outDir, filePath)
	if err != nil {
		videoPathInHtm = filePath
	}

	// windows 的路徑沒辦法直接放進網頁，要改成正斜線
	videoPathInHtm = strings.ReplaceAll(videoPathInHtm, "\\", "/")

	var setlistFilename string
	filename := fileNameWithoutExt(filePath) + ".txt"
	_, err = os.Stat(filename)
	if err == nil {
		setlistFilename = filename
	}

	var videoType string
	switch filepath.Ext(filePath) {
	case ".mp4": videoType = "video/mp4"
	case ".webm": videoType = "video/webm"
	default: log.Fatalln("not support video type, only support mp4 and webm")
	}
    setlistJsonText, err := convertSetlist2Json(setlistFilename)
	if err != nil {
		return err
	}

	chatJsonInFile, err := os.Open(chatJsonFilename)
	if err != nil {
		return err
	}
	defer chatJsonInFile.Close()
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	var liveChatText []string
	fileScanner := bufio.NewScanner(chatJsonInFile)
	fileScanner.Split(bufio.ScanLines)
	downloader := HttpImgDownloader{}
	for fileScanner.Scan() {
		line := fileScanner.Text()
		jsonText, err := preprocessJson(&downloader, line, outDir)
		if err != nil {
			continue
		}
		liveChatText = append(liveChatText, jsonText)
	}
	htmText := templateHtm
	htmText = strings.ReplaceAll(htmText, "{{video}}", videoPathInHtm)
	htmText = strings.ReplaceAll(htmText, "{{title}}", filepath.Base(filePath))
	htmText = strings.ReplaceAll(htmText, "{{video-type}}", videoType)
	htmText = strings.ReplaceAll(htmText, "{{live-chat-json}}", strings.Join(liveChatText,  "\n"))
	htmText = strings.ReplaceAll(htmText, "{{setlist-json}}", setlistJsonText)

	log.Println("output: " + outputFilename)
	outputFile.WriteString(htmText)
	jsFile, err := os.Create(filepath.Join(outDir, "play-live-chat.js"))
	if err != nil {
		return err
	}
	defer jsFile.Close()
	jsFile.WriteString(playlivechatjs)

	return nil
}
