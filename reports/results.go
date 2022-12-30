package reports

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/countcraicula/hytek"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func ResultSheet(m *hytek.Meet, events []*hytek.Event, opts ...SheetOption) ([]bytes.Buffer, error) {
	var ret []bytes.Buffer
	s := applyOptions(opts)
	eventList := [][]*hytek.Event{events}
	if s.BySession() {
		eventList = s.EventOrder().SplitBySession(events)
	}
	for i, events := range eventList {
		buf, err := resultSheet(m, events, s, i+1)
		if err != nil {
			return nil, err
		}
		ret = append(ret, buf)
	}
	return ret, nil
}

func resultSheet(m *hytek.Meet, events []*hytek.Event, s *SheetOptions, session int) (bytes.Buffer, error) {
	p := pdf.NewMaroto(s.Orientation(), s.Size())
	p.SetAliasNbPages("{nb}")
	p.SetFirstPageNb(1)
	p.SetDefaultFontFamily(consts.Courier)
	p.RegisterHeader(resultHeader(p, m, s, session))
	p.RegisterFooter(resultFooter(p))
	o := s.EventOrder()
	o.Sort(events)
	for _, event := range events {
		if len(event.Entries) == 0 {
			continue
		}
		sort.Slice(event.Entries, func(i, j int) bool {
			if event.Entries[i].Entry.Result.Time == 0 {
				return false
			}
			if event.Entries[j].Entry.Result.Time == 0 {
				return true
			}
			return event.Entries[i].Entry.Result.Time < event.Entries[j].Entry.Result.Time
		})
		resultEventHeader(p, event)
		for i, entry := range event.Entries {
			resultEntry(p, entry, i+1)
		}
	}
	return p.Output()
}

func resultHeader(p pdf.Maroto, m *hytek.Meet, s *SheetOptions, session int) func() {
	return func() {
		p.Row(10, func() {
			p.Col(3, func() {
				p.Text(m.Description)

			})
			p.Col(3, func() {
				p.Text(m.Location, props.Text{Align: consts.Center})
			})
			p.Col(3, func() {
				p.Text(fmt.Sprintf("Session %v", session), props.Text{Align: consts.Center, Style: consts.Bold})
			})
			p.Col(3, func() {
				p.Text(s.SessionTime(session).Format("02/01/2006"), props.Text{Align: consts.Right})
			})
		})
		p.Line(1.0)
		p.Row(10, func() {
			p.Col(12, func() {
				p.Text("Result sheet", props.Text{Align: consts.Center})
			})
		})
		p.Line(1.0)
	}
}

func resultFooter(p pdf.Maroto) func() {

	return func() {
		p.Line(1.0)
		p.Row(10, func() {
			p.Col(12, func() {
				p.Text(fmt.Sprintf("Page %v/{nb}", p.GetCurrentPage()), props.Text{Align: consts.Center})
			})
		})
	}
}

func resultEventHeader(p pdf.Maroto, event *hytek.Event) {
	p.Row(10, func() {
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
	p.Row(6, func() {

		p.Col(1, func() {
			p.Text("Place", props.Text{Style: consts.Bold})
		})
		p.ColSpace(1)
		p.Col(4, func() {
			p.Text("Name", props.Text{Style: consts.Bold})
		})
		p.Col(1, func() {
			p.Text("Age", props.Text{Style: consts.Bold})
		})
		p.Col(2, func() {
			p.Text("Result", props.Text{Align: consts.Right, Style: consts.Bold})
		})
		p.Col(2, func() {
			p.Text("Entry Time", props.Text{Align: consts.Right, Style: consts.Bold})
		})
	})
}

func resultEntry(p pdf.Maroto, entry *hytek.Entry, place int) {

	p.Row(6, func() {
		p.Col(1, func() {
			p.Text(fmt.Sprintf("%v.", place), props.Text{Align: consts.Right})
		})
		p.ColSpace(1)
		p.Col(4, func() {
			p.Text(fmt.Sprintf("%v, %v", entry.Swimmer.LastName, entry.Swimmer.FirstName))
		})
		p.Col(1, func() {
			p.Text(fmt.Sprint(entry.Swimmer.Age))
		})
		p.Col(2, func() {
			if entry.Entry.Result.Time == 0 {
				p.Text("NS", props.Text{Align: consts.Right})
			} else {
				p.Text(fmt.Sprintf("%v", entry.Entry.Result.Time), props.Text{Align: consts.Right})
			}
		})
		p.Col(2, func() {
			p.Text(fmt.Sprintf("%v", entry.Entry.SeedTime1), props.Text{Align: consts.Right})
		})
	})
}
