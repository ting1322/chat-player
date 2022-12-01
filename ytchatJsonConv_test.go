package main

import (
	_ "embed"
	"regexp"
	"strings"
	"testing"
)

//go:embed testdata/chat-message.json
var chatjson string

type FakeDownloader struct {
	urlList []string
}

func (me *FakeDownloader) Download(localPath, url string) string {
	me.urlList = append(me.urlList, url)
	return localPath
}

// 聊天室使用 emoji 的紀錄，需要轉換圖片網址為本地檔案，
// 並且要有 ImgDownloader 曾經下載過得紀錄
func TestEmoji(t *testing.T) {
	option = NewOption()
	fd := FakeDownloader{}
	text, err := preprocessJson(&fd, chatjson, "")
	if err != nil {
		t.Fatal(err)
	}
	if text == chatjson {
		t.Fatal("nothing change")
	}
	want := regexp.MustCompile(`\{"url":"images/emoji-7db03403b139a71c9ebe20a29ff51ca2f4a34e10\.svg"\}`)
	if !want.MatchString(text) {
		t.Fatalf("not found pattern, text: %v, \nwant: %#q\n", text, want)
	}
	if len(fd.urlList) != 1 {
		t.Fatal("download item count")
	}
	if fd.urlList[0] != "https://www.youtube.com/s/gaming/emoji/0f0cae22/emoji_u1f31a.svg" {
		t.Fatal("download item")
	}
}

func TestTimeOffset(t *testing.T) {
	option = NewOption()
	option.TimeOffsetInSec = 30
	fd := FakeDownloader{}
	text, err := preprocessJson(&fd, chatjson, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, `"videoOffsetTimeMsec":"30000"`) {
		t.Fatal("video offset not equal 30 * 1000")
	}
}