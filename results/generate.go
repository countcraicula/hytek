package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/countcraicula/hytek"
	"github.com/jszwec/csvutil"

	result "github.com/countcraicula/hytek/csv"
)

var (
	input          = flag.String("input", "", "")
	output         = flag.String("output", "test.hy3", "")
	resultsFile    = flag.String("results", "", "")
	outputTemplate = flag.Bool("output_template", false, "")
)

func main() {
	flag.Parse()

	in, err := os.Open(*input)
	if err != nil {
		fmt.Println(*input)
		fmt.Println(err)
	}
	file, err := hytek.ParseHY3File(in)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *resultsFile != "" {
		r, err := os.Open(*resultsFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		raw := csv.NewReader(r)
		decoder, err := csvutil.NewDecoder(raw)
		if err != nil {
			fmt.Println(err)
			return
		}
		var res []*result.Result
		if err := decoder.Decode(&res); err != nil {
			fmt.Println(err)
			return
		}

		var results = make(map[string]map[string]*result.Result)
		for _, r := range res {
			s, ok := results[r.ID]
			if !ok {
				s = make(map[string]*result.Result)
				results[r.ID] = s
			}
			s["FIXME"] = r
		}

		file.FileDescriptor.Type = "07"

		for _, team := range file.Teams {
			for _, swimmer := range team.Swimmers {
				for _, entry := range swimmer.IndividualEntries {
					r, ok := results[swimmer.Info1.ID][entry.EventNumber]
					if !ok {
						fmt.Printf("Unknown entry for %v:%v\n", swimmer.Info1.SwimmerIDEvent, entry.EventNumber)
						continue
					}
					entry.Result = &hytek.HY3IndividualEventResults{
						Type:     r.Type,
						Time:     r.Time,
						TimeCode: "S", //No finals
						Splits:   r.Splits(),
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
	hytek.GenerateHY3File(file, out)
}
