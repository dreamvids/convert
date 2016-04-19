package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const (
	FormatWebM = 1
	FormatMp4  = 2

	StatusError      = 1
	StatusConverting = 2
	StatusAvailable  = 3
)

type Conversion struct {
	ID       int `json:"id"`
	VideoID  int `json:"video_id"`
	FormatID int `json:"format_id"`
	StatusID int `json:"status_id"`
}

func NewConversion(videoId, formatId, statusId int) Conversion {
	return Conversion{0, videoId, formatId, statusId}
}

func (c *Conversion) Start() error {
	var format string
	var dst string
	var cmd *exec.Cmd

	src := TempDir + strconv.Itoa(c.VideoID) + ".video"

	switch c.FormatID {
	case FormatWebM:
		format = "webm"
		dst = TempDir + strconv.Itoa(c.VideoID) + ".webm"
		cmd = exec.Command("ffmpeg", "-loglevel", "error", "-i", src, "-c:v", "libvpx", "-c:a", "libvorbis", "-f", format, dst)
		break
	case FormatMp4:
		format = "mp4"
		dst = TempDir + strconv.Itoa(c.VideoID) + ".mp4"
		cmd = exec.Command("ffmpeg", "-loglevel", "error", "-i", src, "-c:v", "libx264", "-c:a", "libmp3lame", "-f", format, dst)
		break
	default:
		return fmt.Errorf("Invalid format")
	}

	f, err := os.Create(dst + ".ffmpeg.log")
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go func() {
		io.Copy(f, stdout)
		f.Close()
		stdout.Close()
	}()

	c.StatusID = StatusConverting
	err = DatabaseUpdateConversion(c)
	if err != nil {
		return err
	}

	go func() {
		err := cmd.Run()
		if err != nil {
			c.StatusID = StatusError
			log.Printf("Failed conversion %d (src: %s, dst: %s): %s\n", c.ID, src, dst, err)
		} else {
			c.StatusID = StatusAvailable
			log.Printf("Finished conversion %d (src: %s, dst: %s)\n", c.ID, src, dst)
		}

		err = DatabaseUpdateConversion(c)
		if err != nil {
			log.Println("Can not finish conversion: database:", err)
			return
		}
	}()

	return nil
}
