package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ting1322/chat-player/pkg/cplayer"
)

var (
	programVersion string = "1.x-dev"
)

func parseCommandline() *cplayer.Option {
	option := cplayer.NewOption()
	var argVersion bool
	flag.StringVar(&option.ChatJson, "chat-json", "", "live chat json file (download by yt-dlp)")
	flag.IntVar(&option.TimeOffsetInSec, "offset", 0, "time offset for live chat (second)")
	flag.StringVar(&option.SetList, "set-list", "", "時間軸 txt 檔")
	flag.StringVar(&option.OutputName, "output", "", "output html file, 不指定就是目前工作目錄跟影片同檔名的htm")
	flag.StringVar(&option.OutDir, "out-dir", "", "輸出目錄，預設是目前工作目錄")
	flag.BoolVar(&option.NoDownloadPic, "no-download-pic", false, "不要把聊天室貼圖抓下來 (每次開網頁使用youtube檔案)")
	flag.BoolVar(&option.SplitRes, "split-res", false, "分離 javascript, css 檔案，預設是嵌在html裡面")
	flag.BoolVar(&argVersion, "version", false, "show program version and exit.")
	flag.Parse()

	if argVersion {
		fmt.Println("chatplayer", programVersion)
		return nil
	}

	if len(flag.Args()) > 0 {
		option.Path = flag.Arg(0)
	}
	return option
}

func main() {
	var option = parseCommandline()
	if option == nil {
		return
	}

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
			err = cplayer.ProcessVideo(option, video)
			if err != nil &&
				err != cplayer.ErrNotFoundJson {
				log.Fatalln(err)
			}
		}

	} else {
		err = cplayer.ProcessVideo(option, option.Path)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
