package cplayer

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ProcessVideo(option *Option, videoFile string) error {
	var chatJsonFilename string = option.ChatJson
	if len(chatJsonFilename) == 0 {
		chatJsonFilename = fileNameWithoutExt(videoFile) + ".live_chat.json"
	}
	_, err := os.Stat(chatJsonFilename)
	if err != nil {
		return ErrNotFoundJson
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
	case ".mp4":
		videoType = "video/mp4"
	case ".webm":
		videoType = "video/webm"
	default:
		log.Fatalln("not support video type, only support mp4 and webm")
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
		jsonText, err := preprocessJson(option, &downloader, line, outDir)
		if err != nil {
			continue
		}
		liveChatText = append(liveChatText, jsonText)
	}
	htmText := TemplateHtm
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
		js = "<script>\n" + Playlivechatjs + "\n</script>"
		css = "<style>\n" + StyleCss + "\n</style>"
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
	jsFile.WriteString(Playlivechatjs)
	cssFile, err := os.Create(filepath.Join(outDir, "style.css"))
	if err != nil {
		return err
	}
	defer cssFile.Close()
	cssFile.WriteString(StyleCss)
	return nil
}
