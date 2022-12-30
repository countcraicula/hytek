package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/countcraicula/hytek"
	"github.com/golang/glog"
	"github.com/jszwec/csvutil"

	result "github.com/countcraicula/hytek/csv"
)

var (
	input          = flag.String("input", "", "")
	hyv            = flag.String("hyv", "", "")
	output         = flag.String("output", "test.hy3", "")
	resultsFile    = flag.String("results", "", "")
	outputTemplate = flag.Bool("output_template", false, "")
)

func main() {
	flag.Parse()

	in, err := os.Open(*input)
	if err != nil {
		glog.Error(*input)
		glog.Error(err)
		return
	}
	file, err := hytek.ParseHY3File(in)
	if err != nil {
		glog.Error(err)
		return
	}

	mIn, err := os.Open(*hyv)
	if err != nil {
		glog.Error(err)
		return
	}
	meet, err := hytek.ParseHyv(mIn)
	if err != nil {
		glog.Error(err)
		return
	}

	lookupEvent := func(distance int, stroke hytek.StrokeCode) *hytek.Event {
		for _, event := range meet.Events {
			if event.Distance == distance && event.Stroke == stroke {
				return event
			}
		}
		glog.Infof("Failed to find event: %v %v", distance, stroke)
		return nil
	}

	if *resultsFile != "" {
		r, err := os.Open(*resultsFile)
		if err != nil {
			glog.Error(err)
			return
		}

		raw := csv.NewReader(r)
		decoder, err := csvutil.NewDecoder(raw)
		if err != nil {
			glog.Error(err)
			return
		}
		var res []*result.Result
		if err := decoder.Decode(&res); err != nil {
			glog.Error(err)
			return
		}

		var results = make(map[string]map[key]*result.Result)
		for _, r := range res {
			s, ok := results[r.ID]
			if !ok {
				s = make(map[key]*result.Result)
				results[r.ID] = s
				glog.Infof("Result: %v, %v, %v, %v", r.LastName, r.FirstName, r.ID, r.Time)
			}
			s[resultKey(r)] = r
		}

		file.FileDescriptor.Type = "07"

		for _, team := range file.Teams {
			for _, swimmer := range team.Swimmers {
				entries := swimmer.IndividualEntries
				lookupEntry := func(distance int, stroke hytek.StrokeCode) *hytek.HY3IndividualEventEntryInfo {
					for _, entry := range entries {
						if entry.Distance == distance && entry.Stroke == stroke {
							return entry
						}
					}
					return nil
				}
				swimmer.IndividualEntries = nil
				for _, result := range results[swimmer.Info1.ID] {
					event := lookupEvent(result.Distance, result.Stroke)
					entry := lookupEntry(result.Distance, result.Stroke)
					swimmer.IndividualEntries = append(swimmer.IndividualEntries, resultToEntry(event, swimmer, result, entry))
				}
			}
		}
	}
	out, err := os.Create(*output)
	if err != nil {
		glog.Error(err)
		return
	}
	hytek.GenerateHY3File(file, out)
}

func resultToEntry(event *hytek.Event, s *hytek.HY3Swimmer, r *result.Result, entry *hytek.HY3IndividualEventEntryInfo) *hytek.HY3IndividualEventEntryInfo {
	abbr := s.Info1.LastName
	if len(abbr) > 5 {
		abbr = abbr[0:5]
	}
	if entry == nil {
		entry = &hytek.HY3IndividualEventEntryInfo{
			Stroke:         r.Stroke,
			Distance:       r.Distance,
			Unknown2:       "NN",
			Unknown3:       "N",
			Gender:         event.Gender,
			AgeLower:       fmt.Sprintf("%03d", event.MinAge),
			AgeUpper:       fmt.Sprintf("%03d", event.MaxAge),
			SwimmerIDEvent: s.Info1.SwimmerIDEvent,
			SwimmerAbbr:    abbr,
			Gender1:        event.Gender,
			Gender2:        event.Gender,
			EventNumber:    event.Number,
		}
	}
	entry.Result =
		&hytek.HY3IndividualEventResults{
			Type:       r.Type,
			Time:       r.Time,
			TimeCode:   "S", //No finals
			LengthUnit: string(hytek.ShortMetres),
			Splits:     r.Splits(),
		}
	return entry
}

type key struct {
	stroke   hytek.StrokeCode
	distance int
}

func resultKey(r *result.Result) key {
	return key{
		stroke:   r.Stroke,
		distance: r.Distance,
	}
}

func entryKey(e *hytek.HY3IndividualEventEntryInfo) key {
	return key{
		stroke:   e.Stroke,
		distance: e.Distance,
	}
}
