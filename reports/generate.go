package main

import (
	"flag"
	"fmt"
	"hytek"
	"hytek/csv"
	"os"
	"sort"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
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
	sessions := make([][]*hytek.Event, 3)
	for _, event := range events {
		k := eventKey{
			stroke:   event.Stroke,
			distance: event.Distance,
		}
		v := eventOrder[k]
		if len(event.Entries) == 0 {
			continue
		}
		sessions[v.session-1] = append(sessions[v.session-1], event)
	}

	for i, events := range sessions {
		dateIndex = i
		session = i + 1
		PsychSheet(pdf.NewMaroto(consts.Portrait, consts.A4), fmt.Sprintf("psychsheet-session-%v.pdf", session), m, events)
		HeatSheet(pdf.NewMaroto(consts.Portrait, consts.A4), fmt.Sprintf("heatsheet-session-%v.pdf", session), m, events)
		LaneSheets(pdf.NewMaroto(consts.Portrait, consts.A4), fmt.Sprintf("lanesheet-session-%v.pdf", session), m, events)
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
