package main

import (
	"flag"
	"log"
)

var (
	workdir string

	hour  int
	limit int
)

func init() {
	flag.StringVar(&workdir, "workdir", ".", "work dir")
	flag.IntVar(&hour, "hour", -1, "transcode hour, default value is to transcode now in batch mode")
	flag.IntVar(&limit, "limit", 1, "limit how many files transcode per run")
}

func main() {
	flag.Parse()

	t := NewTranscoder(workdir)

	if hour == -1 {
		log.Println("Starting transcode in BATCH mode")
		log.Println()
		t.StartTranscode(limit)
	} else {
		log.Println("Starting transcode in continues execution mode")
		log.Println()
		t.RunLoop(hour, limit)
	}
}
