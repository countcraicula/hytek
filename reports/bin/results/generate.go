package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/countcraicula/hytek"
	"github.com/countcraicula/hytek/reports"
	"github.com/golang/glog"
)

var (
	resultsFile = flag.String("results", "", "")
	meetFile    = flag.String("meet", "", "")
)

func main() {
	flag.Parse()
	f, err := os.Open(*resultsFile)
	if err != nil {
		glog.Fatalf("Failed to open results file: %v", err)
	}
	file, err := hytek.ParseHY3File(f)
	if err != nil {
		glog.Fatalf("Failed to parse results file: %v", err)
	}

	f, err = os.Open(*meetFile)
	if err != nil {
		glog.Fatalf("Failed to open meet file: %v", err)
	}
	meet, err := hytek.ParseHyv(f)
	if err != nil {
		glog.Fatalf("Failed to parse meet file: %v", err)
	}
	var opts = []reports.SheetOption{
		reports.SessionTimesOption([]time.Time{
			time.Date(2022, 12, 29, 10, 30, 0, 0, time.Local),
			time.Date(2022, 12, 30, 10, 30, 0, 0, time.Local),
		}),
	}
	if err := hytek.PopulateMeetEntries(meet, file); err != nil {
		glog.Fatalf("Failed to populate meet entries: %v", err)
	}

	resultsBuf, err := reports.ResultSheet(meet, meet.Events, opts...)
	if err != nil {
		glog.Fatalf("Failed to generate results: %v", err)
	}

	for i, v := range resultsBuf {
		if err := os.WriteFile(fmt.Sprintf("results-%d.pdf", i+1), v.Bytes(), 0644); err != nil {
			glog.Fatalf("Failed to write results file: %v", err)
		}
	}

}
