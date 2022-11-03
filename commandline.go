package main

import ("flag")

var optionChatJson string
var optionSetList string
var optionOutputName string
var optionNoDownloadPic bool
var optionPath string = "."
var optionOutDir string
var optionSplitRes bool

func parseCommandline() {
	flag.StringVar(&optionChatJson, "chat-json", "", "live chat json file (download by yt-dlp)")
	flag.StringVar(&optionSetList, "set-list", "", "時間軸 txt 檔")
	flag.StringVar(&optionOutputName, "output", "", "output html file, 不指定就是目前工作目錄跟影片同檔名的htm")
	flag.StringVar(&optionOutDir, "out-dir", "", "輸出目錄，預設是目前工作目錄")
	flag.BoolVar(&optionNoDownloadPic, "no-download-pic", false, "不要把聊天室貼圖抓下來 (每次開網頁使用youtube檔案)")
	flag.BoolVar(&optionSplitRes, "split-res", false, "分離 javascript, css 檔案，預設是嵌在html裡面")
	flag.Parse()

	if len(flag.Args()) > 0 {
		optionPath = flag.Arg(0)
	}
}