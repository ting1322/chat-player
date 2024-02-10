#!/bin/bash -x

set -e

ver=`cat version.txt`-`git rev-parse --short HEAD`

rm -f chatplayer \
   chatplayer.exe \
   chatplayer-linux-x86-64.zip \
   chatplayer-windows-x86-64.zip

go test github.com/ting1322/chat-player/pkg/cplayer

go build -o chatplayer -ldflags "-X main.programVersion=$ver"

zip chatplayer-linux-x86-64.zip chatplayer

GOOS=windows go build -o chatplayer -ldflags "-X main.programVersion=$ver"

zip chatplayer-windows-x86-64.zip chatplayer.exe
