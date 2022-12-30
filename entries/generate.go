package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/countcraicula/hytek"
)

var (
	input  = flag.String("input", "", "")
	output = flag.String("output", "test.hyv", "")
)

var events = []struct {
	finals    bool
	ageGroups []hytek.QualifyingTime
	events    []struct {
		event     hytek.StrokeCode
		distances []int
	}
	mixed bool
}{
	{
		finals: false,
		ageGroups: []hytek.QualifyingTime{
			{
				MinAge: 9,
				MaxAge: 99,
			},
		},
		events: []struct {
			event     hytek.StrokeCode
			distances []int
		}{
			{
				event:     hytek.Freestyle,
				distances: []int{50, 100, 200, 400},
			}, {
				event:     hytek.Backstroke,
				distances: []int{50, 100, 200},
			}, {
				event:     hytek.Breaststroke,
				distances: []int{50, 100, 200},
			}, {
				event:     hytek.Butterfly,
				distances: []int{50, 100, 200},
			}, {
				event:     hytek.Medley,
				distances: []int{100, 200},
			},
		},
		mixed: true,
	},
}

func main() {
	flag.Parse()
	m := hytek.Meet{
		Description:     "Christmas Time Trial",
		StartDate:       time.Date(2022, time.December, 27, 10, 30, 0, 0, time.Local),
		EndDate:         time.Date(2022, time.December, 29, 12, 0, 0, 0, time.Local),
		AgeUpDate:       time.Date(2022, time.December, 31, 0, 0, 0, 0, time.Local),
		CourseCode:      hytek.ShortMetres,
		Location:        "Aura Dundalk",
		SoftwareVendor:  "Hy-Tek Sports Software",
		SoftwareVersion: "8.0De",
		Unknown2:        "CN",
	}

	ageGroup := []hytek.QualifyingTime{
		{
			MinAge: 7,
			MaxAge: 99,
		},
	}

	// Day 1
	// 200 free
	// 25 free
	// 50 free
	// 100 back
	// 50 breast
	// 100 fly
	// 200 breast
	// 200 IM
	// 4x25 free relay
	m.AddEvents("1", hytek.Freestyle, hytek.Mixed, 200, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("2", hytek.Freestyle, hytek.Mixed, 25, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("3", hytek.Freestyle, hytek.Mixed, 50, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("4", hytek.Backstroke, hytek.Mixed, 100, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("5", hytek.Breaststroke, hytek.Mixed, 25, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("6", hytek.Breaststroke, hytek.Mixed, 50, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("7", hytek.Butterfly, hytek.Mixed, 100, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("8", hytek.Breaststroke, hytek.Mixed, 200, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("9", hytek.Medley, hytek.Mixed, 200, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("10", hytek.Freestyle, hytek.Mixed, 25, hytek.Relay, ageGroup, hytek.Finals)

	// Day 2
	// 200 back
	// 100 breast
	// 25 back
	// 50 back
	// 100 free
	// 50 fly
	// 100 IM
	// 400 free
	// 4x25 medley relay
	m.AddEvents("11", hytek.Backstroke, hytek.Mixed, 200, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("12", hytek.Breaststroke, hytek.Mixed, 100, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("13", hytek.Backstroke, hytek.Mixed, 25, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("14", hytek.Backstroke, hytek.Mixed, 50, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("15", hytek.Freestyle, hytek.Mixed, 100, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("16", hytek.Butterfly, hytek.Mixed, 50, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("17", hytek.Medley, hytek.Mixed, 100, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("18", hytek.Freestyle, hytek.Mixed, 400, hytek.Individual, ageGroup, hytek.Finals)
	m.AddEvents("19", hytek.Medley, hytek.Mixed, 25, hytek.Relay, ageGroup, hytek.Finals)
	m.AddEvents("20", hytek.Butterfly, hytek.Mixed, 200, hytek.Individual, ageGroup, hytek.Finals)

	out, err := os.Create(*output)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(out, m.String())
}
