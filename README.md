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

# 補充說明

1. 影片必須是 mp4 或 webm，瀏覽器只能播放這兩種。千萬別用 mkv。
2. 聊天室的貼圖會下載並存在 images 目錄，供離線使用。
