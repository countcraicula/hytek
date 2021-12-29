package reports

import (
	"bytes"
	"fmt"

	"github.com/countcraicula/hytek"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func PsychSheet(m *hytek.Meet, events []*hytek.Event, opts ...SheetOption) ([]bytes.Buffer, error) {
	var ret []bytes.Buffer
	s := applyOptions(opts)
	eventList := [][]*hytek.Event{events}
	if s.BySession() {
		eventList = s.EventOrder().SplitBySession(events)
	}
	for i, events := range eventList {
		buf, err := psychSheet(m, events, s, i+1)
		if err != nil {
			return nil, err
		}
		ret = append(ret, buf)
	}
	return ret, nil
}

func psychSheet(m *hytek.Meet, events []*hytek.Event, s *SheetOptions, session int) (bytes.Buffer, error) {
	p := pdf.NewMaroto(s.Orientation(), s.Size())
	p.SetAliasNbPages("{nb}")
	p.SetFirstPageNb(1)
	p.SetDefaultFontFamily(consts.Courier)
	p.RegisterHeader(psychHeader(p, m, s, session))
	p.RegisterFooter(psychFooter(p))
	for _, event := range events {
		if len(event.Entries) == 0 {
			continue
		}
		psychEventHeader(p, event)
		for i, entry := range event.Entries {
			psychEntry(p, entry, i+1)
		}
	}
	return p.Output()
}

func psychHeader(p pdf.Maroto, m *hytek.Meet, s *SheetOptions, session int) func() {
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
				p.Text("Psych sheet", props.Text{Align: consts.Center})
			})
		})
		p.Line(1.0)
	}
}

func psychFooter(p pdf.Maroto) func() {

	return func() {
		p.Line(1.0)
		p.Row(10, func() {
			p.Col(12, func() {
				p.Text(fmt.Sprintf("Page %v/{nb}", p.GetCurrentPage()), props.Text{Align: consts.Center})
			})
		})
	}
}

func psychEventHeader(p pdf.Maroto, event *hytek.Event) {
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
}

func psychEntry(p pdf.Maroto, entry *hytek.Entry, place int) {

	p.Row(6, func() {
		p.Col(1, func() {
			p.Text(fmt.Sprintf("%v.", place), props.Text{Align: consts.Right})
		})
		p.ColSpace(1)
		p.Col(5, func() {
			p.Text(fmt.Sprintf("%v, %v", entry.Swimmer.LastName, entry.Swimmer.FirstName))
		})
		p.Col(2, func() {
			p.Text(fmt.Sprint(entry.Swimmer.Age))
		})
		p.Col(3, func() {
			p.Text(fmt.Sprintf("%v", entry.Entry.SeedTime1), props.Text{Align: consts.Right})
		})
	})
}
