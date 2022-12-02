#!/bin/bash -x

set -e

rm -f chatplayer \
   chatplayer.exe \
   chatplayer-linux-x86-64.zip \
   chatplayer-windows-x86-64.zip

go test github.com/ting1322/chat-player/pkg/cplayer

go build

GOOS=windows go build

zip chatplayer-linux-x86-64.zip chatplayer

zip chatplayer-windows-x86-64.zip chatplayer.exe
