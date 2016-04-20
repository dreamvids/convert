package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Stream struct {
	Index     int    `json:"index"`
	CodecType string `json:"codec_type"`
	Codec     string `json:"codec_name"`
	Duration  string `json:"duration"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	BitRate   string `json:"bit_rate"`
}

func ProbeVideo(path string) (Stream, error) {
	type ProbeResult struct {
		Streams []Stream
	}

	var r ProbeResult
	var ss Stream

	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", path)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return ss, err
	}

	err = cmd.Start()
	if err != nil {
		return ss, err
	}

	err = json.NewDecoder(stdout).Decode(&r)
	if err != nil {
		return ss, err
	}

	err = cmd.Wait()
	if err != nil {
		return ss, err
	}

	for _, s := range r.Streams {
		if s.CodecType == "video" {
			return s, nil
		}
	}

	return ss, fmt.Errorf("Invalid video file (no video stream)")
}
