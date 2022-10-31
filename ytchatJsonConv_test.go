package main

import (
	_ "embed"
	"regexp"
	"regexp"
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

func TestEmoji(t *testing.T) {
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
