package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/countcraicula/hytek"
)

func main() {
	o1, err := os.Open("../data/test.hyv")
	if err != nil {
		fmt.Println(err)
		return
	}
	meet, err := hytek.ParseHyv(o1)

	if err != nil {
		fmt.Println(err)
		return
	}

	o2, err := os.Open("../data/ADSC-Entries-Christmas Time Trial-29Dec2022-001.HY3")
	if err != nil {
		fmt.Println(err)
		return
	}
	entries, err := hytek.ParseHY3File(o2)
	if err != nil {
		fmt.Println(err)
		return
	}
	hytek.PopulateMeetEntries(meet, entries)

	for _, event := range meet.Events {
		event.AssignHeats(3)
	}

	b, err := json.MarshalIndent(meet, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := ioutil.WriteFile("meet.json", b, 0744); err != nil {
		fmt.Println(err)
		return
	}
}
