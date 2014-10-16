// Package transcode implements routines for transcoding to various kinds of
// receiver.
package transcode

import (
	"io"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"bitbucket.org/anacrolix/dms/ffmpeg"
	. "bitbucket.org/anacrolix/dms/misc"
)

// Invokes an external command and returns a reader from its stdout. The
// command is waited on asynchronously.
func transcodePipe(args []string, stderr io.Writer) (r io.ReadCloser, err error) {
	log.Println("transcode command:", args)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = stderr
	r, err = cmd.StdoutPipe()
	if err != nil {
		return
	}
	err = cmd.Start()
	if err != nil {
		return
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("command %s failed: %s", args, err)
		}
	}()
	return
}

// Return a series of ffmpeg arguments that pick specific codecs for specific
// streams. This requires use of the -map flag.
func streamArgs(s map[string]interface{}) (ret []string) {
	defer func() {
		if len(ret) != 0 {
			ret = append(ret, []string{
				"-map", "0:" + strconv.Itoa(int(s["index"].(float64))),
			}...)
		}
	}()
	switch s["codec_type"] {
	case "video":
		/*
			if s["codec_name"] == "h264" {
				if i, _ := strconv.ParseInt(s["is_avc"], 0, 0); i != 0 {
					return []string{"-vcodec", "copy", "-sameq", "-vbsf", "h264_mp4toannexb"}
				}
			}
		*/
		return []string{"-target", "pal-dvd"}
	case "audio":
		if s["codec_name"] == "dca" {
			return []string{"-acodec", "ac3", "-ab", "224k", "-ac", "2"}
		} else {
			return []string{"-acodec", "copy"}
		}
	case "subtitle":
		return []string{"-scodec", "copy"}
	}
	return
}

// Streams the desired file in the MPEG_PS_PAL DLNA profile.
func Transcode(path string, start, length time.Duration, stderr io.Writer) (r io.ReadCloser, err error) {
	args := []string{
		"ffmpeg",
		"-threads", strconv.FormatInt(int64(runtime.NumCPU()), 10),
		"-async", "1",
		"-ss", FormatDurationSexagesimal(start),
	}
	if length >= 0 {
		args = append(args, []string{
			"-t", FormatDurationSexagesimal(length),
		}...)
	}
	args = append(args, []string{
		"-i", path,
	}...)
	info, err := ffmpeg.Probe(path)
	if err != nil {
		return
	}
	for _, s := range info.Streams {
		args = append(args, streamArgs(s)...)
	}
	args = append(args, []string{"-f", "mpegts", "pipe:"}...)
	return transcodePipe(args, stderr)
}

// Returns a stream of Chromecast supported VP8.
func VP8Transcode(path string, start, length time.Duration, stderr io.Writer) (r io.ReadCloser, err error) {
	args := []string{
		"avconv",
		"-threads", strconv.FormatInt(int64(runtime.NumCPU()), 10),
		"-async", "1",
		"-ss", FormatDurationSexagesimal(start),
	}
	if length > 0 {
		args = append(args, []string{
			"-t", FormatDurationSexagesimal(length),
		}...)
	}
	args = append(args, []string{
		"-i", path,
		// "-deadline", "good",
		// "-c:v", "libvpx", "-crf", "10",
		"-f", "webm",
		"pipe:"}...)
	return transcodePipe(args, stderr)
}

// Returns a stream of Chromecast supported matroska.
func ChromecastTranscode(path string, start, length time.Duration, stderr io.Writer) (r io.ReadCloser, err error) {
	args := []string{
		"ffmpeg",
		"-ss", FormatDurationSexagesimal(start),
		"-i", path,
		"-c:v", "libx264", "-preset", "ultrafast", "-profile:v", "high", "-level", "5.0",
		"-movflags", "+faststart",
	}
	if length > 0 {
		args = append(args, []string{
			"-t", FormatDurationSexagesimal(length),
		}...)
	}
	args = append(args, []string{
		"-f", "matroska",
		"pipe:"}...)
	return transcodePipe(args, stderr)
}