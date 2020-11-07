package main

import "flag"

var (
	workdir string

	hour  int
	limit int
)

func init() {
	flag.StringVar(&workdir, "workdir", ".", "work dir")
	flag.IntVar(&hour, "hour", -1, "transcode hour, default value is to transcode now")
	flag.IntVar(&limit, "limit", 1, "limit how many files transcode per run")
}

func main() {
	flag.Parse()

	t := NewTranscoder(workdir)

	if hour == -1 {
		t.StartTranscode(limit)
	} else {
		t.RunLoop(hour, limit)
	}
}
