package main

import (
	"flag"
	"log"
)

var (
	workdir string

	hour  int
	limit int

	autoDelete bool
)

func init() {
	flag.StringVar(&workdir, "workdir", ".", "work dir")
	flag.IntVar(&hour, "hour", -1, "transcode hour, default value is to transcode now in batch mode")
	flag.IntVar(&limit, "limit", 2, "limit how many files transcode per run")
	flag.BoolVar(&autoDelete, "auto-delete", false, "auto delete original video file")
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
