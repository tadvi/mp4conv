package main

import (
	"flag"
	"log"
)

var (
	workDir string
	outDir  string

	hour  int
	limit int

	batch      bool
	autoDelete bool
)

func init() {
	flag.BoolVar(&batch, "batch", false, "run in batch mode")

	flag.StringVar(&workDir, "workdir", ".", "work directory")
	flag.StringVar(&outDir, "outdir", ".", "output directory")

	flag.IntVar(&limit, "limit", 1, "limit how many files transcode per run")
	flag.BoolVar(&autoDelete, "auto-delete", false, "auto delete original video file")
}

func main() {
	flag.Parse()

	t := NewTranscoder(workDir, outDir)

	if batch {
		log.Println("Starting transcode in BATCH mode")
		log.Println()
		t.StartTranscode(limit)
	} else {
		log.Println("Starting transcode in continues execution mode")
		log.Println()
		t.RunLoop(limit)
	}
}
