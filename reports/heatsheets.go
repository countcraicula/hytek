package reports

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/countcraicula/hytek"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

var dateList = []time.Time{
	time.Date(2021, 12, 28, 10, 30, 0, 0, time.Local),
	time.Date(2021, 12, 30, 10, 30, 0, 0, time.Local),
	time.Date(2022, 01, 02, 10, 30, 0, 0, time.Local),
}

func eventToHeatDuration(e *hytek.Event) time.Duration {
	switch e.Distance {
	case 25:
		return 60 * time.Second
	case 50:
		return 90 * time.Second
	case 100:
		return 120 * time.Second
	case 200:
		return 240 * time.Second
	case 400:
		return 480 * time.Second
	}
	return 10 * time.Minute
}

var startTimeFormat = "3:04pm"

func HeatSheet(m *hytek.Meet, events []*hytek.Event, opts ...SheetOption) ([]bytes.Buffer, error) {
	var ret []bytes.Buffer
	s := applyOptions(opts)
	eventList := [][]*hytek.Event{events}
	if s.BySession() {
		eventList = s.EventOrder().SplitBySession(events)
	}
	for i, events := range eventList {
		buf, err := heatSheet(m, events, s, i+1)
		if err != nil {
			return nil, err
		}
		ret = append(ret, buf)
	}
	return ret, nil
}

func heatSheet(m *hytek.Meet, events []*hytek.Event, s *SheetOptions, session int) (bytes.Buffer, error) {
	p := pdf.NewMaroto(s.Orientation(), s.Size())
	p.SetAliasNbPages("{nb}")
	p.SetFirstPageNb(1)
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
				p.Text(s.SessionTime(session).Format("02/01/2006"), props.Text{Align: consts.Right, Style: consts.Bold})
			})
		})
		p.Line(1.0)
		p.Row(6, func() {
			p.Col(12, func() {
				p.Text("Heat sheet", props.Text{Align: consts.Center, Style: consts.Bold})
			})
		})
		p.Line(1.0)
	})
	p.RegisterFooter(func() {
		p.Line(1.0)
		p.Row(10, func() {
			p.Col(12, func() {
				p.Text(fmt.Sprintf("Page %v/{nb}", p.GetCurrentPage()), props.Text{Align: consts.Center})
			})
		})
	})
	generateHeatList(p, m, events, s, session)
	return p.Output()
}

func generateHeatList(p pdf.Maroto, m *hytek.Meet, events []*hytek.Event, s *SheetOptions, session int) {
	startTime := s.SessionTime(session)
	o := s.EventOrder()
	o.Sort(events)
	heatSessionHeader(p, session)
	for _, event := range events {
		if len(event.Entries) == 0 {
			continue
		}
		sort.Sort(sortByHeatAndLane(event.Entries))
		maybeAddPageBeforeEvent(p, 3)
		heatEventHeader(p, event, startTime)
		heat := 0
		for _, entry := range event.Entries {
			if entry.Entry.Result.Heat != heat {
				heat = entry.Entry.Result.Heat
				maybeAddPageBeforeHeat(p, s.Lanes())
				heatHeader(p, heat, startTime)
				startTime = startTime.Add(eventToHeatDuration(event))
			}
			heatEntry(p, entry.Entry.Result.Lane, entry)
		}
		if o.BreakAfter(event) {
			breakHeader(p)
			startTime = startTime.Add(10 * time.Minute)
		}
	}
}

func heatSessionHeader(p pdf.Maroto, session int) {
	p.Row(6, func() {
		p.Col(12, func() {
			p.Text(fmt.Sprintf("Session %v - %v", session, dateList[session-1].Format("02/01/2006 - 03:04pm")), props.Text{Style: consts.Bold, Align: consts.Center})
		})
	})
}

const heatEventHeaderHeight = 20

func heatEventHeader(p pdf.Maroto, event *hytek.Event, startTime time.Time) {
	p.Row(6, func() {})
	p.Line(1.0)
	p.Row(6, func() {
		p.ColSpace(3)
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

func breakHeader(p pdf.Maroto) {
	p.Line(1.0)
	p.Row(6, func() {
		p.Col(12, func() {
			p.Text("10 minute break", props.Text{Style: consts.Bold, Align: consts.Center})
		})
	})
	p.Line(1.0)
}

const heatHeaderHeight = 6

func heatHeader(p pdf.Maroto, heat int, startTime time.Time) {
	p.Row(6, func() {
		p.Col(6, func() {
			p.Text(
				fmt.Sprintf("Heat %v", heat),
				props.Text{Style: consts.Bold})
		})
		p.Col(6, func() {
			p.Text(fmt.Sprintf("Start time: %v", startTime.Format(startTimeFormat)), props.Text{Align: consts.Right, Style: consts.Bold})
		})
	})
}

const heatEntryHeight = 6

func heatEntry(p pdf.Maroto, lane int, entry *hytek.Entry) {
	p.Row(6, func() {
		p.Col(2, func() {
			p.Text(fmt.Sprint(lane), props.Text{Align: consts.Right})
		})
		p.ColSpace(1)
		p.Col(4, func() {
			p.Text(fmt.Sprintf("%v, %v", entry.Swimmer.LastName, entry.Swimmer.FirstName))
		})
		p.ColSpace(1)
		p.Col(4, func() {
			p.Text(entry.Entry.SeedTime1.String(), props.Text{Align: consts.Right})
		})
	})
}
func heatDistanceFromBottom(p pdf.Maroto) float64 {
	_, h := p.GetPageSize()
	_, _, _, b := p.GetPageMargins()
	o := p.GetCurrentOffset()
	return h - b - 11 - o
}

func maybeAddPageBeforeHeat(p pdf.Maroto, numLanes int) {
	d := int(heatDistanceFromBottom(p)) + 1
	h := heatEntryHeight*(numLanes) + heatHeaderHeight + 10
	if d < h {
		p.AddPage()
	}
}

func maybeAddPageBeforeEvent(p pdf.Maroto, numLanes int) {
	d := int(heatDistanceFromBottom(p)) + 1
	h := heatEntryHeight*(numLanes) + heatHeaderHeight + heatEventHeaderHeight + 10
	if d < h {
		p.AddPage()
	}
}
