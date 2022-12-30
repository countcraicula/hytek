package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/countcraicula/hytek"
	"github.com/countcraicula/hytek/csv"
	"github.com/countcraicula/hytek/reports"
	"github.com/jszwec/csvutil"
)

var (
	hy3      = flag.String("hy3", "", "")
	hyv      = flag.String("hyv", "", "")
	numLanes = flag.Int("num_lanes", 3, "")
)

func main() {
	flag.Parse()

	hy3In, err := os.Open(*hy3)
	if err != nil {
		fmt.Println(err)
	}
	hyvIn, err := os.Open(*hyv)
	if err != nil {
		fmt.Println(err)
	}
	m, err := hytek.ParseHyv(hyvIn)
	if err != nil {
		fmt.Println("Failed to parse HYV")
		fmt.Println(err)
		return
	}
	entries, err := hytek.ParseHY3File(hy3In)
	if err != nil {
		fmt.Println("Failed to parse HY3")
		fmt.Println(err)
		return
	}
	if err := hytek.PopulateMeetEntries(m, entries); err != nil {
		fmt.Println("Failed to populate meet entries")
		fmt.Println(err)
		return
	}
	addMastersEvents(m, entries)
	events := m.Events
	for _, event := range events {
		sort.Sort(event.Entries)
		event.AssignHeats(*numLanes)
	}
	var psychOpts = []reports.SheetOption{
		reports.SessionTimesOption([]time.Time{
			time.Date(2022, 12, 29, 10, 30, 0, 0, time.Local),
			time.Date(2022, 12, 30, 10, 30, 0, 0, time.Local),
		}),
		reports.NumLanesOption(*numLanes),
	}
	var opts []reports.SheetOption
	opts = append(opts, psychOpts...)
	opts = append(opts, reports.BySessionOption(true))

	psychBufs, err := reports.PsychSheet(m, events, psychOpts...)
	if err != nil {
		fmt.Println(err)
	}
	for session, buf := range psychBufs {
		os.WriteFile(fmt.Sprintf("psychsheet-%v.pdf", session+1), buf.Bytes(), 0755)
	}
	heatBufs, err := reports.HeatSheet(m, events, opts...)
	if err != nil {
		fmt.Println(err)
	}
	for session, buf := range heatBufs {
		os.WriteFile(fmt.Sprintf("heatsheet-%v.pdf", session+1), buf.Bytes(), 0755)
	}
	laneBufs, err := reports.LaneSheets(m, events, opts...)
	if err != nil {
		fmt.Println(err)
	}
	for session, buf := range laneBufs {
		os.WriteFile(fmt.Sprintf("lanesheet-%v.pdf", session+1), buf.Bytes(), 0755)
	}
	res := csv.MeetToResults(m)
	out, err := os.Create("results.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := csvutil.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	out.Write(b)
	f, err := os.Create("meet.hy3")
	if err != nil {
		fmt.Println(err)
		return
	}
	blankResults(entries)
	if err := hytek.GenerateHY3File(entries, f); err != nil {
		fmt.Println(err)
		return
	}
}

func blankResults(e *hytek.HY3) {
	for _, t := range e.Teams {
		for _, s := range t.Swimmers {
			for _, r := range s.IndividualEntries {
				r.Result = nil
			}
		}
	}
}

type mastersEntry struct {
	Info1  *hytek.HY3SwimmerInfo1
	Events map[eventKey]hytek.HY3Time
}

type eventKey struct {
	stroke   hytek.StrokeCode
	distance int
}

func addMastersEvents(m *hytek.Meet, entries *hytek.HY3) {
	var swimmers = []*mastersEntry{}
	id := 5000
	team := entries.Teams[0]

	for _, swimmer := range swimmers {
		s := &hytek.HY3Swimmer{
			Info1: swimmer.Info1,
		}
		team.Swimmers = append(team.Swimmers, s)
		for event, t := range swimmer.Events {
			e := lookupEvent(m, event.stroke, event.distance)
			if e == nil {
				continue
			}
			abbr := swimmer.Info1.LastName
			if len(abbr) > 5 {
				abbr = abbr[:5]
			}
			entry := &hytek.HY3IndividualEventEntryInfo{
				Gender:         swimmer.Info1.Gender,
				SwimmerIDEvent: id,
				SwimmerAbbr:    abbr,
				Gender1:        e.Gender,
				Gender2:        e.Gender,
				Distance:       event.distance,
				Stroke:         event.stroke,
				AgeLower:       "07",
				AgeUpper:       "109",
				EventNumber:    e.Number,
				SeedTime1:      t,
				SeedCourse1:    string(hytek.ShortMetres),
			}
			s.IndividualEntries = append(s.IndividualEntries, entry)

			e.Entries = append(e.Entries, &hytek.Entry{
				Swimmer: swimmer.Info1,
				Entry:   entry,
			})

		}
		id++
	}
}

func lookupEvent(m *hytek.Meet, stroke hytek.StrokeCode, distance int) *hytek.Event {
	for _, event := range m.Events {
		if event.Stroke == stroke && event.Distance == distance {
			return event
		}
	}
	return nil
}
