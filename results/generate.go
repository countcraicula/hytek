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

		var results = make(map[string]map[key]*result.Result)
		for _, r := range res {
			s, ok := results[r.ID]
			if !ok {
				s = make(map[key]*result.Result)
				results[r.ID] = s
			}
			s[resultKey(r)] = r
		}

		file.FileDescriptor.Type = "07"

		for _, team := range file.Teams {
			for _, swimmer := range team.Swimmers {
				for _, entry := range swimmer.IndividualEntries {
					r, ok := results[swimmer.Info1.ID][entryKey(entry)]
					if !ok {
						fmt.Printf("Unknown entry for %v,%v:%v\n", swimmer.Info1.LastName, swimmer.Info1.FirstName, entryKey(entry))
						continue
					}
					entry.Unknown2 = "NN"
					entry.Unknown3 = "N"
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

type key struct {
	stroke   hytek.StrokeCode
	distance int
}

func resultKey(r *result.Result) key {
	return key{
		stroke:   r.Stroke,
		distance: r.Distance,
	}
}

func entryKey(e *hytek.HY3IndividualEventEntryInfo) key {
	return key{
		stroke:   e.Stroke,
		distance: e.Distance,
	}
}
