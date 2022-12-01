package main

import ("flag")

type Option struct {
	ChatJson string
	SetList string
	OutputName string
	NoDownloadPic bool
	Path string
	OutDir string
	SplitRes bool
	TimeOffsetInSec int
}

func NewOption() *Option {
	var opt = &Option{}
	opt.Path = "."
	return opt
}

func parseCommandline() *Option {
	option := NewOption()
	flag.StringVar(&option.ChatJson, "chat-json", "", "live chat json file (download by yt-dlp)")
	flag.IntVar(&option.TimeOffsetInSec, "offset", 0, "time offset for live chat (second)")
	flag.StringVar(&option.SetList, "set-list", "", "時間軸 txt 檔")
	flag.StringVar(&option.OutputName, "output", "", "output html file, 不指定就是目前工作目錄跟影片同檔名的htm")
	flag.StringVar(&option.OutDir, "out-dir", "", "輸出目錄，預設是目前工作目錄")
	flag.BoolVar(&option.NoDownloadPic, "no-download-pic", false, "不要把聊天室貼圖抓下來 (每次開網頁使用youtube檔案)")
	flag.BoolVar(&option.SplitRes, "split-res", false, "分離 javascript, css 檔案，預設是嵌在html裡面")
	flag.Parse()

	if len(flag.Args()) > 0 {
		option.Path = flag.Arg(0)
	}
	return option
}