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

聯絡方式: 到 discord 「GuildCQ 傳教士公會」 tag 猫耳大好き
