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

//go:embed style.css
var styleCss string

var option *Option

func main() {
	option = parseCommandline()

	fileInfo, err := os.Stat(option.Path)
	if err != nil {
		log.Fatalln(err)
	}
	if fileInfo.IsDir() {
		if len(option.ChatJson) != 0 ||
			len(option.SetList) != 0 {
			log.Fatalln("not support option in directory mode")
		}
		webms, err := filepath.Glob(option.Path + "/*.webm")
		if err != nil {
			log.Fatalln(err)
		}
		mp4s, err := filepath.Glob(option.Path + "/*.mp4")
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
		err = processVideo(option.Path)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

const not_found_json string = "not found .live_chat.json file"

func processVideo(videoFile string) error {
	var chatJsonFilename string = option.ChatJson
	if len(chatJsonFilename) == 0 {
		chatJsonFilename = fileNameWithoutExt(videoFile) + ".live_chat.json"
	}
	_, err := os.Stat(chatJsonFilename)
	if err != nil {
		return errors.New(not_found_json)
	}
	_, err = os.Stat(videoFile)
	if err != nil {
		return err
	}
	var outputFilename string = option.OutputName
	var outDir string = option.OutDir
	if len(option.OutputName) == 0 {
		outputFilename = videoFile + ".htm"
	}
	if len(outDir) > 0 {
		outputFilename = filepath.Join(outDir, filepath.Base(outputFilename))
	} else {
		outDir = filepath.Dir(outputFilename)
	}
	videoPathInHtm := relPathAsUrl(outDir, videoFile)

	var setlistFilename string
	filename := fileNameWithoutExt(videoFile) + ".txt"
	_, err = os.Stat(filename)
	if err == nil {
		setlistFilename = filename
	}

	var videoType string
	switch filepath.Ext(videoFile) {
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
	htmText = strings.ReplaceAll(htmText, "{{title}}", filepath.Base(videoFile))
	htmText = strings.ReplaceAll(htmText, "{{video-type}}", videoType)
	htmText = strings.ReplaceAll(htmText, "{{live-chat-json}}", strings.Join(liveChatText, "\n"))
	htmText = strings.ReplaceAll(htmText, "{{setlist-json}}", setlistJsonText)
	var js string
	var css string
	if option.SplitRes {
		js = `<script src="play-live-chat.js"></script>`
		css = `<link rel="stylesheet" type="text/css" href="style.css">`
		if err := writeResFile(outDir); err != nil {
			return err
		}
	} else {
		js = "<script>\n" + playlivechatjs + "\n</script>"
		css = "<style>\n" + styleCss + "\n</style>"
	}
	htmText = strings.ReplaceAll(htmText, "{{javascript1}}", js)
	htmText = strings.ReplaceAll(htmText, "{{stylecss1}}", css)

	log.Println("output: " + outputFilename)
	outputFile.WriteString(htmText)
	return nil
}

func writeResFile(outDir string) error {
	jsFile, err := os.Create(filepath.Join(outDir, "play-live-chat.js"))
	if err != nil {
		return err
	}
	defer jsFile.Close()
	jsFile.WriteString(playlivechatjs)
	cssFile, err := os.Create(filepath.Join(outDir, "style.css"))
	if err != nil {
		return err
	}
	defer cssFile.Close()
	cssFile.WriteString(styleCss)
	return nil
}