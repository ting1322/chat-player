# chat-player

yt-dlp 下載的 XXX.live-chat.json 聊天室紀錄檔，可以用這個專案播放。
拉動影片時間軸，同步捲動聊天室。
暫停播放狀態，點聊天室的紀錄兩下，把影片跳到指定時間。

# 使用方式

1. 用 yt-dlp 抓影片，加上參數
   ```
   --write-subs --sub-langs live_chat
   ```
2. 用本專案的 generate-htm.py 產生 htm
   ```
   ./generate-htm.py XXX.webm
   ```
3. 瀏覽器開啟 htm 檔

## 完整 command line
```
➜  chat-player git:(main) ./generate-htm.py --help
usage: generate-htm.py [-h] [-c CHAT_JSON] [-s SET_LIST] [-o OUTPUT] [--no-download-pic] video

generate a html to play video with live-chat.json

positional arguments:
  video                 video file (webm or mp4)

optional arguments:
  -h, --help            show this help message and exit
  -c CHAT_JSON, --chat-json CHAT_JSON
                        live chat json file (download by yt-dlp)
  -s SET_LIST, --set-list SET_LIST
                        時間軸 txt 檔
  -o OUTPUT, --output OUTPUT
                        output html file
  --no-download-pic     不要把聊天室貼圖抓下來 (每次開網頁使用youtube檔案)
```

只有輸入的影片檔路徑必須要給，其餘皆是可選。
輸出檔案預設在當前工作目錄，檔名為影片檔 + .htm。
json 檔案預設是影片檔名 + live_chat.json，這也是 yt-dlp 下載下來預設的檔名。

# 安裝環境

解壓縮得到 generate-htm.py, play-live-chat.js, template.htm.in 就是全部檔案。
電腦必須先安裝 python 3.9 以上版本。

# 補充說明

1. 影片必須是 mp4 或 webm，瀏覽器只能播放這兩種。千萬別用 mkv。
2. 聊天室的貼圖會下載並存在 images 目錄，供離線使用。
   可以加參數 --no-download-pic 關閉這個下載功能，讓網頁使用線上的圖片。
   （不建議，考慮會員貼圖有刪除的可能）
3. 貼圖存檔的檔名，採用原始網址hash的值。相同貼圖、相同解析度只會下載一次。
   images 目錄可以供多個 htm 共用。
4. 可以提供時間軸 txt 檔，格式範例如下
   ````
   0:00:00 start
   0:01:30 chapter 1
   0:25:00 chapter 2
   ````
   文字檔如果放在 XXX.webm 旁邊，命名為 set-list.txt，會自動讀入。
   其他檔案名稱可以用 --set-list FILENAME.txt 輸入
5. 我的環境是 Linux Ubuntu 21.10, python 3.9, firefox 97.0。
   反應問題請附上軟體環境。
6. 測試過 久遠たま、苺咲べりぃ 的影片。如果其他人的影片有問題，給我 json 檔看看。
   
# implement detail

底下列出相同目的，失敗的嘗試。如果你打算做相同的東西，這些經驗或許能節省點時間。

- 將 live-chat.json 轉檔為 ass 字幕，跟著影片播放。
  遇到問題: 貼圖無法顯示。libass 並不支援 picture event，雖然spec有寫。
- 用 html + javascript 載入本機 json
  瀏覽器安全問題，擋掉開啟本機另一個檔案的功能。
- 用 SimpleHTTPServer 解決前一個問題
  有 http server 確實能讓 javascript 順利載入 json。但是影片播放無法 seek。

所以目前實做方式: 所有檔案都嵌入同一個 htm 檔案，讓瀏覽器開啟本機檔案。
影片 seek 沒問題，文字轉 json 沒問題。

聯絡方式: 到 discord 「GuildCQ 傳教士公會」 tag 猫耳大好き
