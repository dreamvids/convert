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
	var format string
	var resolution string
	var bitrate string
	var dst string
	var cmd *exec.Cmd

	src := TempDir + strconv.Itoa(c.VideoID) + ".video"

	switch c.ResolutionID {
	case Resolution360p:
		resolution = "640x360"
		bitrate = "1000k"
		break
	case Resolution720p:
		resolution = "1280x720"
		bitrate = "5000k"
		break
	case Resolution1080p:
		resolution = "1920x1080"
		bitrate = "8000k"
		break
	default:
		return fmt.Errorf("Invalid resolution")
	}

	switch c.FormatID {
	case FormatWebM:
		format = "webm"
		dst = TempDir + strconv.Itoa(c.VideoID) + "." + resolution + ".webm"
		cmd = exec.Command("ffmpeg", "-loglevel", "warning", "-i", src, "-b:v", bitrate, "-b:a", "128k", "-c:v", "libvpx", "-c:a", "libvorbis", "-f", format, "-s", resolution, "-framerate", "30", dst)
		break
	case FormatMp4:
		format = "mp4"
		dst = TempDir + strconv.Itoa(c.VideoID) + "." + resolution + ".mp4"
		cmd = exec.Command("ffmpeg", "-loglevel", "warning", "-i", src, "-b:v", bitrate, "-b:a", "128k", "-c:v", "libx264", "-c:a", "libmp3lame", "-f", format, "-s", resolution, "-framerate", "30", dst)
		break
	default:
		return fmt.Errorf("Invalid format")
	}

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
