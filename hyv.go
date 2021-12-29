package hytek

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

type EventClassification string

const (
	Prelims EventClassification = "P"
	Finals  EventClassification = "F"
)

type Gender string

const (
	Male   Gender = "M"
	Female Gender = "F"
	Mixed  Gender = "X"
)

func (g Gender) Display() string {
	return map[Gender]string{
		Male:   "Boys",
		Female: "Girls",
		Mixed:  "Mixed",
	}[g]
}

type EventType string

const (
	Individual EventType = "I"
	Relay      EventType = "R"
)

type StrokeCode int

const (
	Freestyle StrokeCode = iota + 1
	Backstroke
	Breaststroke
	Butterfly
	Medley
)

var strokeCodeMapping = []string{
	"Freestyle",
	"Backstroke",
	"Breaststroke",
	"Butterfly",
	"Medley",
}

func (s StrokeCode) Display() string {
	return strokeCodeMapping[int(s)-1]
}
func (s StrokeCode) MarshalCSV() ([]byte, error) {
	return []byte(strokeCodeMapping[int(s)-1]), nil
}
func (s *StrokeCode) UnmarshalCSV(b []byte) error {
	str := string(b)
	for i, stroke := range strokeCodeMapping {
		if str == stroke {
			*s = StrokeCode(i + 1)
		}
	}
	return nil
}

type Entry struct {
	Swimmer    *HY3SwimmerInfo1
	Entry      *HY3IndividualEventEntryInfo
	RelayEntry *HY3RelayEventEntryInfo
}

type Entries []*Entry

func (e Entries) Less(i, j int) bool {
	if e[i].Entry.SeedTime1 == 0 && e[j].Entry.SeedTime1 != 0 {
		return false
	}
	if e[j].Entry.SeedTime1 == 0 && e[i].Entry.SeedTime1 != 0 {
		return true
	}
	if e[i].Entry.SeedTime1 == e[j].Entry.SeedTime1 {
		return e[i].Swimmer.Age > e[j].Swimmer.Age
	}
	return e[i].Entry.SeedTime1 < e[j].Entry.SeedTime1
}
func (e Entries) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e Entries) Len() int      { return len(e) }

type Event struct {
	Number         string
	Classification EventClassification
	Gender         Gender
	Type           EventType
	MinAge         int
	MaxAge         int
	Distance       int
	Stroke         StrokeCode
	Unknown1       string
	QualifyingTime time.Duration
	Unknown2       string
	EventFee       float32
	Unknown3       string
	Unknown4       string
	Unknown5       string
	ConversionTime time.Duration
	Unknown6       string
	Unknown7       string
	Entries        Entries
}

func sortByLaneOrder(h []*Entry) {
	switch len(h) {
	case 0, 1:
	case 2:
		h[0], h[1] = h[1], h[0]
	case 3:
		h[1], h[2] = h[2], h[1]
	case 4:
		h[0], h[1], h[3] = h[1], h[3], h[0]
	}
}

func (e *Event) AssignHeats(numLanes int) {

	numEntries := len(e.Entries)
	if numEntries == 0 {
		return
	}
	numHeats := len(e.Entries) / numLanes
	modHeats := len(e.Entries) % numLanes
	if modHeats != 0 {
		numHeats++
	}
	currEntry := numEntries - 1
	for i := 0; i < numHeats; i++ {
		numSwimmers := numLanes
		heat := i + 1
		if i == 0 {
			if modHeats != 0 {
				numSwimmers = modHeats
			}
			if modHeats == 1 && numEntries > 1 {
				numSwimmers++
			}
		}
		if i == 1 {
			if modHeats == 1 {
				numSwimmers--
			}
		}
		var heatSwimmers []*Entry
		for j := 0; j < numSwimmers; j++ {
			heatSwimmers = append(heatSwimmers, e.Entries[currEntry])
			currEntry--
		}
		sortByLaneOrder(heatSwimmers)
		for j, entry := range heatSwimmers {
			lane := j + 1
			if numSwimmers < 3 {
				lane++
			}
			entry.Entry.Result = &HY3IndividualEventResults{
				Lane: lane,
				Heat: heat,
			}
		}
	}

}

func (e *Event) String() string {
	return fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v",
		e.Number,
		e.Classification,
		e.Gender,
		e.Type,
		e.MinAge,
		e.MaxAge,
		e.Distance,
		e.Stroke,
		e.Unknown1,
		timeString(e.QualifyingTime),
		e.Unknown2,
		e.EventFee,
		e.Unknown3,
		e.Unknown4,
		e.Unknown5,
		timeString(e.ConversionTime),
		e.Unknown6,
		e.Unknown7,
	)
}

func (e *Event) MarshalText() ([]byte, error) {

	return []byte(fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v",
		e.Number,
		e.Classification,
		e.Gender,
		e.Type,
		e.MinAge,
		e.MaxAge,
		e.Distance,
		e.Stroke,
		e.Unknown1,
		timeString(e.QualifyingTime),
		e.Unknown2,
		e.EventFee,
		e.Unknown3,
		e.Unknown4,
		e.Unknown5,
		timeString(e.ConversionTime),
		e.Unknown6,
		e.Unknown7,
	)), nil
}

func (e *Event) UnmarshalText(b []byte) error {
	var err error
	s := string(b)
	ss := strings.Split(s, ";")
	if len(ss) < 17 {
		return fmt.Errorf("string to short")
	}
	e.Number = ss[0]
	e.Classification = EventClassification(ss[1])
	e.Gender = Gender(ss[2])
	e.Type = EventType(ss[3])
	fmt.Sscan(ss[4], &e.MinAge)
	fmt.Sscan(ss[5], &e.MaxAge)
	fmt.Sscan(ss[6], &e.Distance)
	fmt.Sscan(ss[7], &e.Stroke)
	e.Unknown1 = ss[8]
	if ss[9] != "" {
		e.QualifyingTime, err = time.ParseDuration(ss[9])
		if err != nil {
			return fmt.Errorf("failed to parse Qualifying time %q: %v", ss[9], err)
		}
	}
	e.Unknown2 = ss[10]
	fmt.Sscan(ss[11], &e.EventFee)
	e.Unknown3 = ss[12]
	e.Unknown4 = ss[13]
	e.Unknown5 = ss[14]
	if ss[15] != "" {
		e.ConversionTime, err = time.ParseDuration(ss[15])
		if err != nil {
			return fmt.Errorf("failed to parse ConversionTime time %q: %v", ss[15], err)
		}
	}
	e.Unknown6 = ss[16]
	e.Unknown7 = ss[17]
	return nil
}

func timeString(t time.Duration) string {
	if t == 0 {
		return ""
	}
	return fmt.Sprint(t.Truncate(time.Millisecond * 10))
}

type CourseCode string

const (
	ShortMetres CourseCode = "S"
	ShortYards  CourseCode = "SY"
	LongMeters  CourseCode = "L"
)

const (
	dateFormat = "01/02/2006"
)

type Meet struct {
	Description     string
	StartDate       time.Time
	EndDate         time.Time
	AgeUpDate       time.Time
	CourseCode      CourseCode
	Location        string
	Unknown1        string
	SoftwareVendor  string
	SoftwareVersion string
	Unknown2        string
	Events          []*Event
}

func (m *Meet) String() string {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, m.header())
	for _, v := range m.Events {
		fmt.Fprintln(&buf, v.String())
	}
	return buf.String()
}

func (m *Meet) header() string {
	s := fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%v;%v;%v",
		m.Description,
		m.StartDate.Format(dateFormat),
		m.EndDate.Format(dateFormat),
		m.AgeUpDate.Format(dateFormat),
		m.CourseCode,
		m.Location,
		m.Unknown1,
		m.SoftwareVendor,
		m.SoftwareVersion,
		m.Unknown2)
	chk := generateChecksum(s)
	return fmt.Sprintf("%v;%v", s, chk)
}

func (m *Meet) MarshalText() ([]byte, error) {
	s := fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%v;%v;%v",
		m.Description,
		m.StartDate.Format(dateFormat),
		m.EndDate.Format(dateFormat),
		m.AgeUpDate.Format(dateFormat),
		m.CourseCode,
		m.Location,
		m.Unknown1,
		m.SoftwareVendor,
		m.SoftwareVersion,
		m.Unknown2)
	chk := generateChecksum(s)
	return []byte(fmt.Sprintf("%v;%v", s, chk)), nil
}

func (m *Meet) UnmarshalText(b []byte) error {
	s := string(b)
	ss := strings.Split(s, ";")
	if len(ss) < 10 {
		return fmt.Errorf("string too short")
	}
	m.Description = ss[0]
	start, err := time.Parse(dateFormat, ss[1])
	if err != nil {
		return err
	}
	m.StartDate = start
	end, err := time.Parse(dateFormat, ss[2])
	if err != nil {
		return err
	}
	m.EndDate = end
	age, err := time.Parse(dateFormat, ss[3])
	if err != nil {
		return err
	}
	m.AgeUpDate = age
	m.CourseCode = CourseCode(ss[4])
	m.Location = ss[5]
	m.Unknown1 = ss[6]
	m.SoftwareVendor = ss[7]
	m.SoftwareVersion = ss[8]
	m.Unknown2 = ss[9]
	return nil
}

func generateChecksum(s string) string {
	ss := strings.Split(s, ";")
	sum := 0
	for _, v := range ss {
		for _, w := range v {
			sum += int(w)
		}
	}
	sum /= 7
	sum += 205
	c := fmt.Sprintf("%4.4v", sum)
	return c[3:4] + c[0:3] + s[2:3]
}

type QualifyingTime struct {
	MinAge         int
	MaxAge         int
	QualifyingTime time.Duration
	ConversionTime time.Duration
}

func (m *Meet) AddEvents(eventNumber string, s StrokeCode, g Gender, distance int, t EventType, q []QualifyingTime, c EventClassification) {
	for i, v := range q {
		maxage := v.MaxAge
		if maxage == 0 || maxage < v.MinAge {
			maxage = 109
		}
		event := &Event{
			Stroke:         s,
			Gender:         g,
			Distance:       distance,
			Type:           t,
			QualifyingTime: v.QualifyingTime,
			ConversionTime: v.ConversionTime,
			MinAge:         v.MinAge,
			MaxAge:         maxage,
			Classification: c,
		}
		if len(q) == 1 {
			event.Number = eventNumber
		} else {
			event.Number = fmt.Sprintf("%v%v", eventNumber, string([]rune{rune('A' + i)}))
		}
		m.Events = append(m.Events, event)
	}
}

func ParseHyv(r io.Reader) (*Meet, error) {
	scanner := bufio.NewScanner(r)
	m := &Meet{}
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read first line")
	}
	if err := m.UnmarshalText(scanner.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to parse first line: %v", err)
	}
	line := 1
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) == 0 {
			continue
		}
		e := &Event{}
		if err := e.UnmarshalText(b); err != nil {
			return nil, fmt.Errorf("failed to parse line %v: %v", line, err)
		}
		m.Events = append(m.Events, e)
		line++
	}
	return m, nil
}
