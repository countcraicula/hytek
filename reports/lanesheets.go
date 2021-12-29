package main

import (
	"fmt"
	"hytek"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func collectByLaneNumber(events []*hytek.Event) [][]*hytek.Event {
	ret := make([][]*hytek.Event, *numLanes)
	for _, event := range events {
		tmp := make([][]*hytek.Entry, *numLanes)
		if len(event.Entries) == 0 {
			continue
		}
		for _, entry := range event.Entries {
			tmp[entry.Entry.Result.Lane-1] = append(tmp[entry.Entry.Result.Lane-1], entry)
		}
		for lane, entries := range tmp {
			e := *event
			e.Entries = entries
			ret[lane] = append(ret[lane], &e)
		}
	}
	return ret
}

func LaneSheets(p pdf.Maroto, filename string, m *hytek.Meet, events []*hytek.Event) {
	eventsByLane := collectByLaneNumber(events)
	currLane := 1
	p.SetDefaultFontFamily(consts.Courier)
	p.RegisterHeader(func() {
		p.Row(6, func() {
			p.Col(3, func() {
				p.Text(m.Description, props.Text{Style: consts.Bold})

			})
			p.Col(3, func() {
				p.Text(m.Location, props.Text{Align: consts.Center, Style: consts.Bold})
			})
			p.Col(3, func() {
				p.Text(fmt.Sprintf("Session %v", session), props.Text{Align: consts.Center, Style: consts.Bold})
			})
			p.Col(3, func() {
				p.Text(m.StartDate.Format("02/01/2006"), props.Text{Align: consts.Right, Style: consts.Bold})
			})
		})
		p.Line(1.0)
		p.Row(6, func() {
			p.Col(12, func() {
				p.Text(fmt.Sprintf("Lane sheet - Lane %v", currLane), props.Text{Align: consts.Center, Style: consts.Bold})
			})
		})
		p.Line(1.0)
	})
	p.RegisterFooter(func() {
		p.Line(1.0)
		p.Row(10, func() {
		})
	})
	for lane, events := range eventsByLane {
		lane++
		currLane = lane
		for _, event := range events {
			if len(event.Entries) == 0 {
				continue
			}
			maybeLaneAddPageBeforeEvent(p, event.Entries)
			heat := 1
			laneEventHeader(p, event)
			if len(event.Entries) == 0 {
				laneEventEntry(p, heat, lane, nil)
			}
			for _, entry := range event.Entries {
				for entry.Entry.Result.Heat != heat {
					laneEventEntry(p, heat, lane, nil)
					heat++
				}
				laneEventEntry(p, heat, lane, entry)
				heat++
			}
		}
		p.AddPage()
	}
	p.OutputFileAndClose(filename)
}

const laneFooterHeight = 11
const laneEventHeaderHeight = 13
const laneEventEntryHeight = 10

func laneDistanceFromBottom(p pdf.Maroto) float64 {
	_, h := p.GetPageSize()
	_, _, _, b := p.GetPageMargins()
	o := p.GetCurrentOffset()
	return h - b - laneFooterHeight - o
}

func laneEventHeader(p pdf.Maroto, event *hytek.Event) {
	p.Row(6, func() {
		p.Col(3, func() {
			p.Text(
				fmt.Sprintf("Mixed %v+", event.MinAge),
				props.Text{Style: consts.Bold})
		})
		p.Col(3, func() {
			p.Text(
				fmt.Sprintf("%vm %v", event.Distance, event.Stroke.Display()),
				props.Text{Style: consts.Bold})
		})
	})
	p.Line(1.0)
	p.Row(6, func() {})
}

func laneEventEntry(p pdf.Maroto, heat, lane int, entry *hytek.Entry) {
	p.Row(10, func() {
		p.Col(3, func() {
			p.Text(fmt.Sprintf("Heat %v, Lane %v", heat, lane))
		})
		if entry != nil {
			p.Col(5, func() {
				p.Text(fmt.Sprintf("%v, %v", entry.Swimmer.LastName, entry.Swimmer.FirstName))
			})
		} else {
			p.Col(5, func() {
				p.Text("_____________, __________")
			})
		}
		p.Col(4, func() {
			p.Text("_______  _______  _______")
		})
	})
}

func maybeLaneAddPageBeforeEvent(p pdf.Maroto, entries []*hytek.Entry) {
	last := entries[len(entries)-1]
	numHeats := last.Entry.Result.Heat
	d := int(laneDistanceFromBottom(p)) + 1
	h := laneFooterHeight + laneEventHeaderHeight + laneEventEntryHeight*numHeats
	if d < h {
		p.AddPage()
	}
}
