package main

import "hytek"

type eventKey struct {
	stroke   hytek.StrokeCode
	distance int
}
type eventValue struct {
	order      int
	breakAfter bool
	session    int
}

var eventOrder = map[eventKey]eventValue{
	{stroke: hytek.Freestyle, distance: 100}:    {order: 1, session: 1},
	{stroke: hytek.Breaststroke, distance: 50}:  {order: 2, session: 1, breakAfter: true},
	{stroke: hytek.Backstroke, distance: 50}:    {order: 3, session: 1},
	{stroke: hytek.Butterfly, distance: 100}:    {order: 4, session: 1, breakAfter: true},
	{stroke: hytek.Freestyle, distance: 400}:    {order: 5, session: 1},
	{stroke: hytek.Breaststroke, distance: 100}: {order: 6, session: 2},
	{stroke: hytek.Freestyle, distance: 50}:     {order: 7, session: 2, breakAfter: true},
	{stroke: hytek.Butterfly, distance: 50}:     {order: 8, session: 2},
	{stroke: hytek.Breaststroke, distance: 200}: {order: 9, session: 2},
	{stroke: hytek.Backstroke, distance: 100}:   {order: 10, session: 2, breakAfter: true},
	{stroke: hytek.Medley, distance: 200}:       {order: 11, session: 2},
	{stroke: hytek.Medley, distance: 100}:       {order: 12, session: 2},
	{stroke: hytek.Backstroke, distance: 200}:   {order: 13, session: 3, breakAfter: true},
	{stroke: hytek.Freestyle, distance: 200}:    {order: 14, session: 3, breakAfter: true},
	{stroke: hytek.Medley, distance: 400}:       {order: 15, session: 3},
}

type sortEventsByStrokeAndDistance []*hytek.Event

func (a sortEventsByStrokeAndDistance) Len() int      { return len(a) }
func (a sortEventsByStrokeAndDistance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortEventsByStrokeAndDistance) Less(i, j int) bool {
	iKey := eventKey{
		stroke:   a[i].Stroke,
		distance: a[i].Distance,
	}
	jKey := eventKey{
		stroke:   a[j].Stroke,
		distance: a[j].Distance,
	}
	return eventOrder[iKey].order < eventOrder[jKey].order
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
