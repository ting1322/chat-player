# [下載單獨.exe](https://github.com/ting1322/chat-player/releases/latest/download/generate-htm.exe)

也可以下載 tar.gz 之後，執行裡面的 .py。

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
3. 把程式 (.py 或是 .exe) 放在影片資料夾的上層，然後執行他。
   如果有檔案
   ```
   D:\vtuber\maisaki-berry\2022-04-07-メン限\2022-04-07.mp4
   D:\vtuber\maisaki-berry\2022-04-07-メン限\2022-04-07.live_chat.json
   ```
   那就把程式放在 D:\vtuber\maisaki-berry\generate-htm.py，然後執行他。
   **注意** 只能使用 mp4 與 webm。mkv 沒辦法播放。
4. 會看到產生htm檔案，還有.js檔案與 image 目錄，用 firefox 開 htm 檔。

如果遇到錯誤，說找不到 requests 之類的，先打指令安裝 requests，再回去前一步驟。
(如果是使用.exe版本的話這個步驟幫不上忙)
   ```
   py -m pip install requests
   ```

## 完整 command line

```
$ ./generate-htm.py --help
usage: generate-htm.py [-h] [-c CHAT_JSON] [-s SET_LIST] [-o OUTPUT] [--no-download-pic] path

generate a html to play video with live-chat.json

positional arguments:
  path                  video file (webm or mp4), or directory (find *.webm and *.json recursive)

optional arguments:
  -h, --help            show this help message and exit
  -c CHAT_JSON, --chat-json CHAT_JSON
                        live chat json file (download by yt-dlp)
  -s SET_LIST, --set-list SET_LIST
                        時間軸 txt 檔
  -o OUTPUT, --output OUTPUT
                        output html file, 不指定就是目前工作目錄跟影片同檔名的htm
  --no-download-pic     不要把聊天室貼圖抓下來 (每次開網頁使用youtube檔案)
 ```

只有輸入的影片檔路徑必須要給，其餘皆是可選。
輸出檔案預設在當前工作目錄，檔名為影片檔 + .htm。
json 檔案預設是影片檔名 + live_chat.json，這也是 yt-dlp 下載下來預設的檔名。

第一個參數可以使用 . 代表當前路徑，或 .. 代表上一層路徑。
使用資料夾作為參數時，會去搜尋所有子資料夾，相同檔名的 *.webm 與 *.live\_chat.json。

# 執行環境

我這邊的環境是

- Python 3.9
- Firefox 99.0

如果有問題，先檢查 python 版本，然後是網頁瀏覽器。

解壓縮得到 generate-htm.py, play-live-chat.js, template.htm.in 就是全部檔案。

- generate-htm.py: 用來產生 htm 檔，並下載emoji貼圖。
- play-live-chat.js: htm 網頁開啟時，載入 json 聊天室內容，並同步時間軸。
- template.htm.in: 產生 htm 的 template。

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
5. 測試過 久遠たま、伊冬ユナ、苺咲べりぃ的影片。如果其他人的影片有問題，給我 json 檔看看。
   有貓耳的話，處理速度會比較快。
   
# implement detail

底下列出相同目的，失敗的嘗試。如果你打算做相同的東西，這些經驗或許能節省點時間。

- 將 live-chat.json 轉檔為 ass 字幕，跟著影片播放。
  遇到問題: 貼圖無法顯示。libass 並不支援 picture event，雖然spec有寫。
- 用 html + javascript 載入本機 json。
  瀏覽器安全問題，擋掉開啟本機另一個檔案的功能。
- 用 SimpleHTTPServer 解決前一個問題。
  有 http server 確實能讓 javascript 順利載入 json。但是影片播放無法 seek。

所以目前實做方式: 所有檔案都嵌入同一個 htm 檔案，讓瀏覽器開啟本機檔案。
影片 seek 沒問題，文字轉 json 沒問題。
