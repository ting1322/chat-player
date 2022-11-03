# [下載單獨.exe](https://github.com/ting1322/chat-player/releases/latest/download/chatplayer.exe)

yt-dlp 下載的 XXX.live-chat.json 聊天室紀錄檔，可以用這個專案播放。

# 功能

- 下載聊天室的貼圖，供離線使用。
- 滑鼠拉動影片時間軸時，同步捲動聊天室。
- 點聊天室的紀錄的時間，把影片跳到指定時間。
- 時間軸 (特定格式.txt)

# 使用方式

1. 你要先學會用 [yt-dlp](https://github.com/yt-dlp/yt-dlp) ，抓影片加上參數抓聊天室 json 檔案
   ```
   --write-subs --sub-langs live_chat
   ```
2. 安裝 python 3 (如果是使用.exe版本的話就不用)
3. 把程式放在影片的旁邊，然後執行他。
   如果有檔案
   ```
   D:\vtuber\maisaki-berry\2022-04-07-メン限\2022-04-07.mp4
   D:\vtuber\maisaki-berry\2022-04-07-メン限\2022-04-07.live_chat.json
   ```
   那就把程式放在 D:\vtuber\maisaki-berry\2022-04-07-メン限\chatplayer.exe，然後執行他。
   **注意** 只能使用 mp4 與 webm。mkv 沒辦法播放。
4. 會看到產生htm檔案，還有.js檔案與 image 目錄，用 firefox 開 htm 檔。

## 完整 command line

```
chatplay [option] [file or directory]

option:
  -chat-json string
    	live chat json file (download by yt-dlp)
  -no-download-pic
    	不要把聊天室貼圖抓下來 (每次開網頁使用youtube檔案)
  -out-dir string
    	輸出目錄，預設是目前工作目錄
  -output string
    	output html file, 不指定就是目前工作目錄跟影片同檔名的htm
  -set-list string
    	時間軸 txt 檔
  -split-res
        分離 javascript, css 檔案，預設是嵌在html裡面
 ```

輸出檔案預設在影片檔旁邊，檔名為影片檔 + .htm。
json 檔案預設是影片檔名 + live_chat.json，這也是 yt-dlp 下載下來預設的檔名。

# 補充說明

1. 影片必須是 mp4 或 webm，瀏覽器只能播放這兩種。千萬別用 mkv。
2. 聊天室的貼圖會下載並存在 images 目錄，供離線使用。
   可以加參數 --no-download-pic 關閉這個下載功能，讓網頁使用線上的圖片。
   （不建議，考慮會員貼圖有刪除的可能）
3. 貼圖存檔的檔名，採用原始網址hash的值。相同貼圖、相同解析度只會下載一次。
   images 目錄可以供多個 htm 共用。
4. 可以提供時間軸 txt 檔，格式範例如下
   ```
   0:00:00 start
   0:01:30 chapter 1
   0:25:00 chapter 2
   ```
   文字檔如果放在 XXX.webm 旁邊，命名為 set-list.txt，會自動讀入。
   其他檔案名稱可以用 --set-list FILENAME.txt 輸入。
5. 一個比較無關的事情，直播中、預期直播結束會砍檔的影片想要備份聊天室，
   你需要同時開兩隻程式，可以選擇 yt-dlp + yt-dlp 開兩次，或是 yt-dlp + ytarchive。
   以前者來說，第一個 yt-dlp 需要下 `--write-subs --sub-langs live_chat --no-download`，
   專門下載聊天室，並開第二個 yt-dlp 加參數 `--no-write-subs` 不要聊天室只下載影片。
   第二個 yt-dlp 可以換成 ytarchive。
6. 測試過 久遠たま、伊冬ユナ、苺咲べりぃ的影片。如果其他人的影片有問題，給我 json 檔看看。
   有貓耳的話，處理速度會比較快。
   
# change log

從新到舊

## 2022-11-04 v0.5.0
1. play-live-chat.js 預設嵌入 htm 之中，並新增參數 -split-res 讓檔案獨立出來
2. 修正 windows 產生的影片檔案路徑錯誤
3. 修正 webp 圖片副檔名錯誤

## 2022-10-31 v0.4.0

1. 支援顯示 super sticker、member free message、membership gift
2. 自動縮放影片寬度，當視窗不是 1920 放到最大時，舊版顯示不太對勁
3. 預設的 .htm 輸出目錄改為影片的旁邊 (可以用 -out-dir 指定)
4. 轉檔程式改用 golang 取代 python

## 2022-05-22 v0.3.3

1. 支援目錄作為參數，抓目錄下所有 webm 與 mp4
2. 修正 windows 10 的 edge 瀏覽器按時間沒辦法跳到影片時間

## 2022-05-03 v0.3.1

第一個能在 windows 跑的版本 (先前都有問題，而我沒試過)
