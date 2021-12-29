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
	output = flag.String("output", "test.hy3", "")
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
				MaxAge: 10,
			},
			{
				MinAge: 11,
				MaxAge: 12,
			},
			{
				MinAge: 13,
				MaxAge: 14,
			},
			{
				MinAge: 15,
				MaxAge: 109,
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
				distances: []int{100, 200, 400},
			},
		},
		mixed: false,
	},
}

func main() {
	flag.Parse()
	m := hytek.Meet{
		Description:     "Christmas Time Trial",
		StartDate:       time.Date(2021, time.December, 28, 10, 30, 0, 0, time.Local),
		EndDate:         time.Date(2021, time.December, 30, 12, 0, 0, 0, time.Local),
		AgeUpDate:       time.Date(2021, time.December, 31, 0, 0, 0, 0, time.Local),
		CourseCode:      hytek.ShortMetres,
		Location:        "Aura Dundalk",
		SoftwareVendor:  "Hy-Tek Sports Software",
		SoftwareVersion: "8.0De",
		Unknown2:        "CN",
	}
	count := 1
	for _, v := range events {
		if !v.finals {
			continue
		}
		for _, event := range v.events {
			for _, distance := range event.distances {
				if v.mixed {
					m.AddEvents(fmt.Sprint(count), event.event, hytek.Mixed, distance, hytek.Individual, v.ageGroups, hytek.Prelims)
					count++
				} else {
					for _, gender := range []hytek.Gender{hytek.Female, hytek.Male} {
						m.AddEvents(fmt.Sprint(count), event.event, gender, distance, hytek.Individual, v.ageGroups, hytek.Prelims)
						count++
					}
				}
			}
		}
	}
	for _, v := range events {
		for _, event := range v.events {
			for _, distance := range event.distances {
				if v.mixed {
					m.AddEvents(fmt.Sprint(count), event.event, hytek.Mixed, distance, hytek.Individual, v.ageGroups, hytek.Finals)
					count++
				} else {
					for _, gender := range []hytek.Gender{hytek.Female, hytek.Male} {
						m.AddEvents(fmt.Sprint(count), event.event, gender, distance, hytek.Individual, v.ageGroups, hytek.Finals)
						count++
					}
				}
			}
		}
	}

	out, err := os.Create(*output)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(out, m.String())
}
