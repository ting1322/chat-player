module chatplayer

go 1.20

replace github.com/ting1322/chat-player/pkg/cplayer => ./pkg/cplayer

require github.com/ting1322/chat-player/pkg/cplayer v0.0.0-00010101000000-000000000000

require (
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
)
