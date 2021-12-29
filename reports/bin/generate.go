package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	numLanes = flag.Int("num_lanes", 4, "")
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
		fmt.Println(err)
		return
	}
	events := mergeEvents(m)
	for _, event := range events {
		event.AssignHeats(*numLanes)
	}

	var opts = []reports.SheetOption{
		reports.SessionTimesOption([]time.Time{
			time.Date(2021, 12, 28, 10, 30, 0, 0, time.Local),
			time.Date(2021, 12, 30, 10, 30, 0, 0, time.Local),
			time.Date(2022, 01, 02, 10, 30, 0, 0, time.Local),
		}),
		reports.BySessionOption(true),
		reports.NumLanesOption(*numLanes),
	}

	psychBufs, err := reports.PsychSheet(m, events, opts...)
	if err != nil {
		fmt.Println(err)
	}
	for session, buf := range psychBufs {
		ioutil.WriteFile(fmt.Sprintf("psychsheet-%v.pdf", session+1), buf.Bytes(), 0755)
	}
	heatBufs, err := reports.HeatSheet(m, events, opts...)
	if err != nil {
		fmt.Println(err)
	}
	for session, buf := range heatBufs {
		ioutil.WriteFile(fmt.Sprintf("heatsheet-%v.pdf", session+1), buf.Bytes(), 0755)
	}
	laneBufs, err := reports.LaneSheets(m, events, opts...)
	if err != nil {
		fmt.Println(err)
	}
	for session, buf := range laneBufs {
		ioutil.WriteFile(fmt.Sprintf("lanesheet-%v.pdf", session+1), buf.Bytes(), 0755)
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
}

func mergeEvents(m *hytek.Meet) []*hytek.Event {
	var ret []*hytek.Event
	type key struct {
		distance int
		stroke   hytek.StrokeCode
	}
	tmp := make(map[key]int)
	for _, event := range m.Events {
		k := key{
			distance: event.Distance,
			stroke:   event.Stroke,
		}
		i, ok := tmp[k]
		if !ok {
			e := &hytek.Event{}
			b, _ := event.MarshalText()
			e.UnmarshalText(b)
			ret = append(ret, e)
			i = len(ret) - 1
			tmp[k] = i
		}
		ret[i].Entries = append(ret[i].Entries, event.Entries...)
		sort.Sort(ret[i].Entries)
	}
	return ret
}
