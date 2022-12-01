package cplayer

type Option struct {
	ChatJson        string
	SetList         string
	OutputName      string
	NoDownloadPic   bool
	Path            string
	OutDir          string
	SplitRes        bool
	TimeOffsetInSec int
}

func NewOption() *Option {
	var opt = &Option{}
	opt.Path = "."
	return opt
}
