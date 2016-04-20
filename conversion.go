package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const (
	FormatWebM = 1
	FormatMp4  = 2

	Resolution360p  = 1
	Resolution720p  = 2
	Resolution1080p = 3

	StatusError      = 1
	StatusConverting = 2
	StatusAvailable  = 3
)

type Conversion struct {
	ID           int `json:"id"`
	VideoID      int `json:"video_id"`
	FormatID     int `json:"format_id"`
	ResolutionID int `json:"resolution_id"`
	StatusID     int `json:"status_id"`
}

func NewConversion(videoId, formatId, resolutionId, statusId int) *Conversion {
	return &Conversion{0, videoId, formatId, resolutionId, statusId}
}

func (c *Conversion) Start() error {
	var resolution string
	var format string
	var dst string
	var cmd *exec.Cmd

	src := TempDir + strconv.Itoa(c.VideoID) + ".video"
	opts := make([]string, 19)
	opts[0] = "-loglevel"
	opts[1] = "warning"
	opts[2] = "-i"
	opts[3] = src
	opts[4] = "-b:a"
	opts[5] = "128k"
	opts[6] = "-framerate"
	opts[7] = "30"
	opts[8] = "-f"

	switch c.FormatID {
	case FormatWebM:
		format = "webm"
		opts[9] = "webm"
		opts[10] = "-c:a"
		opts[11] = "libvorbis"
		opts[12] = "-c:v"
		opts[13] = "libvpx"
		break
	case FormatMp4:
		format = "mp4"
		opts[9] = "mp4"
		opts[10] = "-c:a"
		opts[11] = "libmp3lame"
		opts[12] = "-c:v"
		opts[13] = "libx264"
		break
	default:
		return fmt.Errorf("Invalid format")
	}

	switch c.ResolutionID {
	case Resolution360p:
		resolution = "360p"
		opts[14] = "-b:v"
		opts[15] = "1000k"
		opts[16] = "-s"
		opts[17] = "640x360"
		break
	case Resolution720p:
		resolution = "720p"
		opts[14] = "-b:v"
		opts[15] = "5000k"
		opts[16] = "-s"
		opts[17] = "1280x720"
		break
	case Resolution1080p:
		resolution = "1080p"
		opts[14] = "-b:v"
		opts[15] = "8000k"
		opts[16] = "-s"
		opts[17] = "1920x1080"
		break
	default:
		return fmt.Errorf("Invalid resolution")
	}

	dst = TempDir + strconv.Itoa(c.VideoID) + "." + resolution + "." + format
	opts[18] = dst
	cmd = exec.Command("ffmpeg", opts...)

	f, err := os.Create(dst + ".ffmpeg")
	if err != nil {
		return err
	}

	cmd.Stdout = f
	cmd.Stderr = f

	c.StatusID = StatusConverting
	err = DatabaseUpdateConversion(c)
	if err != nil {
		return err
	}

	go func() {
		defer f.Close()

		log.Printf("Starting conversion %d (src: %s, dst: %s)\n", c.ID, src, dst)

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
