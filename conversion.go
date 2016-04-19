package main

const (
	FormatWebM = 1
	FormatMp4  = 2

	StatusError      = 1
	StatusConverting = 2
	StatusAvailable  = 3
)

type Conversion struct {
	ID       int
	VideoID  int
	FormatID int
	StatusID int
}

func NewConversion(videoId, formatId, statusId int) Conversion {
	return Conversion{0, videoId, formatId, statusId}
}
