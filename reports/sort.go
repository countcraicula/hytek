package reports

import (
	"sort"

	"github.com/countcraicula/hytek"
)

type eventKey struct {
	stroke   hytek.StrokeCode
	distance int
}
type eventValue struct {
	order      int
	breakAfter bool
	session    int
}

var DefaultEventOrder = []OrderFunc{
	MixedGenderStrokeDistanceOrder(hytek.Freestyle, 100),
	MixedGenderStrokeDistanceOrder(hytek.Breaststroke, 50),
	BreakOrder(),
	MixedGenderStrokeDistanceOrder(hytek.Backstroke, 50),
	MixedGenderStrokeDistanceOrder(hytek.Butterfly, 100),
	BreakOrder(),
	MixedGenderStrokeDistanceOrder(hytek.Freestyle, 400),
	NewSessionOrder(),
	MixedGenderStrokeDistanceOrder(hytek.Breaststroke, 100),
	MixedGenderStrokeDistanceOrder(hytek.Freestyle, 50),
	BreakOrder(),
	MixedGenderStrokeDistanceOrder(hytek.Butterfly, 50),
	MixedGenderStrokeDistanceOrder(hytek.Breaststroke, 200),
	MixedGenderStrokeDistanceOrder(hytek.Backstroke, 100),
	BreakOrder(),
	MixedGenderStrokeDistanceOrder(hytek.Medley, 100),
	MixedGenderStrokeDistanceOrder(hytek.Medley, 200),
	NewSessionOrder(),
	MixedGenderStrokeDistanceOrder(hytek.Backstroke, 200),
	BreakOrder(),
	MixedGenderStrokeDistanceOrder(hytek.Freestyle, 200),
	BreakOrder(),
	MixedGenderStrokeDistanceOrder(hytek.Medley, 400),
}

type sortEventsByStrokeAndDistance struct {
	events []*hytek.Event
	order  map[eventKey]*eventValue
}

func (a *sortEventsByStrokeAndDistance) Len() int { return len(a.events) }
func (a *sortEventsByStrokeAndDistance) Swap(i, j int) {
	a.events[i], a.events[j] = a.events[j], a.events[i]
}
func (a *sortEventsByStrokeAndDistance) Less(i, j int) bool {
	iKey := eventKey{
		stroke:   a.events[i].Stroke,
		distance: a.events[i].Distance,
	}
	jKey := eventKey{
		stroke:   a.events[j].Stroke,
		distance: a.events[j].Distance,
	}
	return a.order[iKey].order < a.order[jKey].order
}

type sortByHeatAndLane []*hytek.Entry

func (a sortByHeatAndLane) Len() int      { return len(a) }
func (a sortByHeatAndLane) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortByHeatAndLane) Less(i, j int) bool {
	if a[i].Entry.Result.Heat == a[j].Entry.Result.Heat {
		return a[i].Entry.Result.Lane < a[j].Entry.Result.Lane
	}
	return a[i].Entry.Result.Heat < a[j].Entry.Result.Heat
}

type Order struct {
	data    map[eventKey]*eventValue
	last    *eventValue
	session int
	event   int
}

func (o *Order) setEvent(k eventKey, v *eventValue) {
	o.data[k] = v
	o.last = v
}

func (o *Order) Sort(e []*hytek.Event) {
	sort.Sort(&sortEventsByStrokeAndDistance{
		events: e,
		order:  o.data,
	})

}

func (o *Order) BreakAfter(e *hytek.Event) bool {
	key := eventKey{
		stroke:   e.Stroke,
		distance: e.Distance,
	}
	v, ok := o.data[key]
	if !ok {
		return false
	}
	return v.breakAfter
}

func (o *Order) SplitBySession(events []*hytek.Event) [][]*hytek.Event {
	sessions := make([][]*hytek.Event, o.session)
	for _, event := range events {
		k := eventKey{
			stroke:   event.Stroke,
			distance: event.Distance,
		}
		v := o.data[k]
		if len(event.Entries) == 0 {
			continue
		}
		sessions[v.session-1] = append(sessions[v.session-1], event)
	}
	return sessions
}

func NewOrder(order ...OrderFunc) *Order {
	o := &Order{
		data:    make(map[eventKey]*eventValue),
		session: 1,
	}
	for _, f := range order {
		if f(o) {
			o.event++
		}
	}
	return o
}

type OrderFunc func(*Order) bool

func MixedGenderStrokeDistanceOrder(stroke hytek.StrokeCode, distance int) OrderFunc {
	return OrderFunc(func(o *Order) bool {
		key := eventKey{
			stroke:   stroke,
			distance: distance,
		}
		o.setEvent(key, &eventValue{order: o.event, session: o.session})
		return true
	})
}

func BreakOrder() OrderFunc {
	return OrderFunc(func(o *Order) bool {
		if o.last == nil {
			return false
		}
		o.last.breakAfter = true
		return false
	})
}

func NewSessionOrder() OrderFunc {
	return OrderFunc(func(o *Order) bool {
		o.session++
		return false
	})
}
