package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var webTranscodeExtensions = map[string]bool{
	".avi":  true,
	".mkv":  true,
	".ts":   true,
	".mpeg": true,
	".mpg":  true,
	".m4a":  true,
}

type Transcoder struct {
	workdir        string
	alreadyTouched map[string]bool

	mu sync.Mutex
}

func NewTranscoder(workdir string) *Transcoder {
	return &Transcoder{workdir: workdir, alreadyTouched: make(map[string]bool)}
}

func (t *Transcoder) RunLoop(hour, limit int) {
	for {
		now := time.Now()
		if now.Hour() == hour && now.Minute() == 0 {
			t.StartTranscode(limit)
		}
		time.Sleep(time.Minute)
	}
}

func (t *Transcoder) StartTranscode(limit int) {
	// prevent running more than one transcode at the same time
	t.mu.Lock()
	defer t.mu.Unlock()

	var files []string
	err := filepath.Walk(t.workdir, func(path string, info os.FileInfo, err error) error {

		ext := filepath.Ext(path)
		if webTranscodeExtensions[ext] {
			fmt.Printf("Adding file: %s\n", path)
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Printf("Failed to read files for transcode: %v\n", err)
	}

	var i int

	for _, path := range files {
		if _, ok := t.alreadyTouched[path]; ok {
			continue
		}

		if i >= limit {
			break
		}

		log.Println(" --- ")
		log.Println("Transcode start for:", path)
		t.alreadyTouched[path] = true

		if err := t.transcode(path); err != nil {
			i++
			log.Println("Transcode finished for:", path)
		}
	}
}

func (t *Transcoder) filenames(srcname string) (string, string, string) {
	srcname = filepath.Clean(srcname)
	dir := filepath.Dir(srcname)           // "/some dir"
	ext := filepath.Ext(srcname)           // ".avi"
	base := filepath.Base(srcname)         // "somewhere.avi"
	noext := strings.TrimSuffix(base, ext) // "somewhere"

	tmpname := fmt.Sprintf("%s/.%s.mp4", dir, noext)
	dstname := fmt.Sprintf("%s/%s.mp4", dir, noext)
	return srcname, tmpname, dstname
}

func (t *Transcoder) transcode(srcname string) error {
	srcname, tmpname, dstname := t.filenames(srcname)

	if _, err := os.Stat(dstname); os.IsExist(err) {
		log.Printf("Destination file exists %q skipping\n", dstname)
		return err
	}

	srcfi, err := os.Stat(srcname)
	if err != nil {
		log.Printf("Error: job %q: %v\n", srcname, err)
		return err
	}

	// Find ffmpeg
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Printf("Error: can not find ffmpeg: %v\n", err)
		return err
	}

	cmd := exec.Command(ffmpeg,
		"-y",
		"-i", srcname,
		"-vcodec", "libx264",
		"-acodec", "aac",
		"-movflags", "faststart", // make streaming work
		"-preset", "veryfast",
		tmpname,
	)

	/*cmd := exec.Command(ffmpeg,
		"-y",
		"-i", srcname,
		"-codec:v", "libx264",
		"-crf", "25",
		"-bf", "2",
		"-flags", "+cgop",
		"-pix_fmt", "yuv420p",
		"-codec:a", "aac",
		"-strict", "-2",
		"-b:a", "384k",
		"-r:a", "48000",
		"-movflags", "faststart", // make streaming work
		"-max_muxing_queue_size", "500", // handle sparse audio/video frames (see: https://trac.ffmpeg.org/ticket/6375#comment:2)
		tmpname,
	)*/

	// Add as a running job.
	log.Printf("Starting transcode job %q -> %q\n", srcname, dstname)

	// Transcode
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error: job %q:\n%s\n", srcname, string(output))
		// Remove the temp file if it still exists at this point.
		os.Remove(tmpname)
		return err
	}
	log.Printf("Success: job %q:\n%s\n", srcname, string(output))

	// Rename temp file to real file.
	if err := os.Rename(tmpname, dstname); err != nil {
		log.Printf("Error: job %q: %s\n", srcname, err)
		return err
	}

	// check that our new file is a reasonable size.
	// TODO: ffprobe and check duration matches?
	minsize := srcfi.Size() / 5
	dstfi, err := os.Stat(dstname)
	if err != nil {
		log.Printf("Error: job %q: %s\n", srcname, err)
		return err
	}
	if dstfi.Size() < minsize {
		log.Printf("Error: job %q: transcoded is too small (%d vs %d); deleting.\n", srcname, dstfi.Size(), minsize)
		if err := os.Remove(dstname); err != nil {
			log.Println(err)
		}
		return err
	}

	// Remove the source file.
	if autoDelete {
		if err := os.Remove(srcname); err != nil {
			log.Printf("Error: job %q: %s\n", srcname, err)
			return nil
		}
	}
	return nil
}
