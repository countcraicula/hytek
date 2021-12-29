package csv

import (
	"encoding/csv"
	"io"

	"github.com/countcraicula/hytek"
	"github.com/jszwec/csvutil"
)

type Result struct {
	ID            string           `csv:"ID"`
	Race          string           `csv:"Race"`
	LastName      string           `csv:"LastName"`
	FirstName     string           `csv:"FirstName"`
	Stroke        hytek.StrokeCode `csv:"Stroke"`
	Distance      int              `csv:"Distance"`
	Type          hytek.EventClassification
	Time          hytek.HY3Time     `csv:"Time"`
	TimeCode      hytek.HY3TimeCode `csv:"Code"`
	Split1        hytek.HY3Time     `csv:"Split1"`
	Split2        hytek.HY3Time     `csv:"Split2"`
	Split3        hytek.HY3Time     `csv:"Split3"`
	Split4        hytek.HY3Time     `csv:"Split4"`
	Split5        hytek.HY3Time     `csv:"Split5"`
	Split6        hytek.HY3Time     `csv:"Split6"`
	Split7        hytek.HY3Time     `csv:"Split7"`
	Split8        hytek.HY3Time     `csv:"Split8"`
	DQDescription string            `csv:"DQ description"`
	DQCode        string            `csv:"DQ code"`
}

func MeetToResults(m *hytek.Meet) Results {
	var ret Results
	for _, event := range m.Events {
		if len(event.Entries) == 0 {
			continue
		}
		for _, entry := range event.Entries {
			r := &Result{
				ID:        entry.Swimmer.ID,
				LastName:  entry.Swimmer.LastName,
				FirstName: entry.Swimmer.FirstName,
				Stroke:    event.Stroke,
				Distance:  event.Distance,
				Type:      event.Classification,
				Time:      entry.Entry.Result.Time,
				TimeCode:  entry.Entry.Result.TimeCode,
			}
			if entry.Entry.Result.DQDescription != nil {
				r.DQDescription = entry.Entry.Result.DQDescription.Description
				r.DQCode = entry.Entry.Result.DQDescription.Code
			}
			ret = append(ret, r)
		}
	}
	return ret
}

func (r *Result) Splits() (res []*hytek.HY3Splits) {
	if r.Split1 == 0 {
		return
	}
	split := &hytek.HY3Splits{
		Times: make(hytek.HY3SplitTimes, 0),
	}
	res = append(res, split)
	split.Times = append(split.Times, &hytek.HY3SplitTime{
		Length: 2,
		Time:   r.Split1,
	})
	if r.Split2 == 0 {
		return
	}
	split.Times = append(split.Times, &hytek.HY3SplitTime{
		Length: 4,
		Time:   r.Split2,
	})
	if r.Split3 == 0 {
		return
	}
	split.Times = append(split.Times, &hytek.HY3SplitTime{
		Length: 6,
		Time:   r.Split3,
	})
	if r.Split4 == 0 {
		return
	}
	split.Times = append(split.Times, &hytek.HY3SplitTime{
		Length: 8,
		Time:   r.Split4,
	})
	if r.Split5 == 0 {
		return
	}
	split.Times = append(split.Times, &hytek.HY3SplitTime{
		Length: 10,
		Time:   r.Split5,
	})
	if r.Split6 == 0 {
		return
	}
	split.Times = append(split.Times, &hytek.HY3SplitTime{
		Length: 12,
		Time:   r.Split6,
	})
	if r.Split7 == 0 {
		return
	}
	split.Times = append(split.Times, &hytek.HY3SplitTime{
		Length: 14,
		Time:   r.Split7,
	})
	if r.Split8 == 0 {
		return
	}
	split.Times = append(split.Times, &hytek.HY3SplitTime{
		Length: 16,
		Time:   r.Split8,
	})
	return
}

type Results []*Result

func (r Results) Write(w io.Writer) error {
	e := csvutil.NewEncoder(csv.NewWriter(w))
	if err := e.EncodeHeader(&Result{}); err != nil {
		return err
	}
	for _, result := range r {
		if err := e.Encode(result); err != nil {
			return err
		}
	}
	return nil
}

func (r Results) Parse(rr io.Reader) error {
	d, err := csvutil.NewDecoder(csv.NewReader(rr))
	if err != nil {
		return err
	}
	if err := d.Decode(&r); err != nil {
		return err
	}
	return nil
}
