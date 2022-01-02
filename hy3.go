package hytek

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"

	fixedwidth "github.com/countcraicula/go-fixedwidth"
)

type HY3Line string

func (h *HY3Line) setType(s string) { *h = HY3Line(s) }

type typeSetter interface {
	setType(string)
}

type HY3 struct {
	FileDescriptor *HY3FileDescriptor
	MeetInfo       *HY3MeetInfo
	MeetAddress    *HY3MeetAddress
	MeetContact    *HY3MeetContact
	Teams          []*HY3SwimTeam
}

func PopulateMeetEntries(m *Meet, h *HY3) error {
	events := make(map[string]*Event)
	for _, event := range m.Events {
		events[event.Number] = event
	}
	for _, team := range h.Teams {
		for _, swimmer := range team.Swimmers {
			for _, entry := range swimmer.IndividualEntries {
				e, ok := events[entry.EventNumber]
				if !ok {
					return fmt.Errorf("unknown event number %v", entry.EventNumber)
				}
				e.Entries = append(e.Entries, &Entry{Swimmer: swimmer.Info1, Entry: entry})
			}
		}
	}
	for _, e := range m.Events {
		sort.Sort(e.Entries)
	}
	return nil
}

type HY3FileDescriptor struct {
	HY3Line         `fixed:"1,2"`
	Type            string `fixed:"3,4"`
	TypeDescription string `fixed:"5,29"`
	VendorName      string `fixed:"30,44"`
	SoftwareVersion string `fixed:"45,58"`
	Date            string `fixed:"59,66"`
	Time            string `fixed:"68,75,right"`
	LicencedTo      string `fixed:"76,128"`
}

type HY3MeetInfo struct {
	HY3Line   `fixed:"1,2"`
	Name      string `fixed:"3,47"`
	Facility  string `fixed:"48,92"`
	Start     string `fixed:"93,100"`
	End       string `fixed:"101,108"`
	AgeUp     string `fixed:"109,116"`
	Elevation string `fixed:"117,121,right"`
}

type HY3MeetAddress struct {
	HY3Line  `fixed:"1,2"`
	Unknown1 string     `fixed:"3,94"`
	Masters  string     `fixed:"95,96"`
	Type     string     `fixed:"97,98"`
	Course   CourseCode `fixed:"99,99"`
	Unknown2 string     `fixed:"100,100"`
	Unknown3 float32    `fixed:"101,106,right"`
	Course2  CourseCode `fixed:"107,107"`
}

type HY3MeetContact struct {
	HY3Line `fixed:"1,2"`
	Unknown string `fixed:"3,128"`
}

type HY3SwimTeam struct {
	Name     *HY3SwimTeamNameInfo
	Address  *HY3SwimTeamAddressInfo
	Contact  *HY3SwimTeamContactInfo
	Swimmers []*HY3Swimmer
}

type HY3SwimTeamNameInfo struct {
	HY3Line   `fixed:"1,2"`
	Abbr      string `fixed:"3,7"`
	Name      string `fixed:"8,37"`
	ShortName string `fixed:"38,53"`
	LSC       string `fixed:"54,55"`
	Contact1  string `fixed:"56,85"`
	Contact2  string `fixed:"86,105"`
	Unknown1  string `fixed:"119,119"`
	Type      string `fixed:"120,122,right"`
}

type HY3SwimTeamAddressInfo struct {
	HY3Line      `fixed:"1,2"`
	MailTo       string `fixed:"3,32"`
	Address      string `fixed:"33,62"`
	City         string `fixed:"63,92"`
	State        string `fixed:"93,94"`
	ZIP          string `fixed:"95,104"`
	Country      string `fixed:"105,107"`
	Unknown      string `fixed:"108,108"`
	Registration string `fixed:"109,112"`
}

type HY3SwimTeamContactInfo struct {
	HY3Line      `fixed:"1,2"`
	Unknown      string `fixed:"3,32"`
	DaytimePhone string `fixed:"33,52"`
	EveningPhone string `fixed:"53,72"`
	Fax          string `fixed:"73,92"`
	Email        string `fixed:"93,128"`
}

type HY3Swimmer struct {
	Info1             *HY3SwimmerInfo1
	Info2             *HY3SwimmerInfo2
	Info3             *HY3SwimmerInfo3
	Info4             *HY3SwimmerInfo4
	Info5             *HY3SwimmerInfo5
	IndividualEntries []*HY3IndividualEventEntryInfo
}

type HY3SwimmerInfo1 struct {
	HY3Line        `fixed:"1,2"`
	Gender         Gender `fixed:"3,3"`
	SwimmerIDEvent int    `fixed:"4,8,right"`
	LastName       string `fixed:"9,28"`
	FirstName      string `fixed:"29,48"`
	NickName       string `fixed:"49,68"`
	MiddleInitial  string `fixed:"69,69"`
	ID             string `fixed:"70,77,right"`
	SwimmerIDTeam  int    `fixed:"84,88,right"`
	Birth          string `fixed:"89,96"`
	Age            int    `fixed:"98,99"`
	Unknown1       int    `fixed:"105,105"`
	Unknown2       string `fixed:"113,115,right"`
	N              string `fixed:"125,125"`

	USSNumGenerated string `fixed:"70,83"`
}

type HY3SwimmerInfo2 struct {
	HY3Line `fixed:"1,2"`
}

type HY3SwimmerInfo3 struct {
	HY3Line `fixed:"1,2"`
}

type HY3SwimmerInfo4 struct {
	HY3Line `fixed:"1,2"`
}

type HY3SwimmerInfo5 struct {
	HY3Line `fixed:"1,2"`
}

type HY3IndividualEventEntryInfo struct {
	HY3Line             `fixed:"1,2"`
	Gender              Gender         `fixed:"3,3"`
	SwimmerIDEvent      int            `fixed:"4,8,right"`
	SwimmerAbbr         string         `fixed:"9,13"`
	Gender1             Gender         `fixed:"14,14"`
	Gender2             Gender         `fixed:"15,15"`
	Distance            int            `fixed:"18,21,right"`
	Stroke              StrokeCode     `fixed:"22,22"`
	AgeLower            string         `fixed:"23,25,right"`
	AgeUpper            string         `fixed:"26,28,right"`
	Unknown1            string         `fixed:"29,32,right"`
	EventFee            float32        `fixed:"33,38,right"`
	EventNumber         string         `fixed:"39,42,right"`
	ConversionSeedTime1 HY3DefaultTime `fixed:"43,50,right"`
	ConversionCourse1   string         `fixed:"51,51"`
	SeedTime1           HY3Time        `fixed:"53,59,right"`
	SeedCourse1         string         `fixed:"60,60"`
	ConversionSeedTime2 HY3Time        `fixed:"61,68,right"`
	ConversionCourse2   string         `fixed:"69,69"`
	SeedTime2           HY3Time        `fixed:"70,76,right"`
	SeedCourse2         string         `fixed:"77,77"`
	Unknown2            string         `fixed:"80,81"`
	Unknown3            string         `fixed:"97,97"`
	Result              *HY3IndividualEventResults
}

type HY3IndividualEventResults struct {
	HY3Line       `fixed:"1,2"`
	Type          EventClassification `fixed:"3,3"`
	Time          HY3Time             `fixed:"4,11,right"`
	LengthUnit    string              `fixed:"12,12"`
	TimeCode      HY3TimeCode         `fixed:"13,15,right"`
	Unknown1      string              `fixed:"16,20,right"`
	Heat          int                 `fixed:"21,23,right"`
	Lane          int                 `fixed:"24,26,right"`
	PlaceInHeat   int                 `fixed:"27,29,right"`
	PlaceOverall  int                 `fixed:"30,33,right"`
	Unknown2      int                 `fixed:"34,36,right"`
	Time1         HY3PlungerTime      `fixed:"37,44,right"`
	Time2         HY3PlungerTime      `fixed:"45,52,right"`
	Time3         HY3PlungerTime      `fixed:"53,60,right"`
	Time4         HY3PlungerTime      `fixed:"66,73,right"`
	Time5         HY3PlungerTime      `fixed:"75,82,right"`
	ReactionTime  HY3ReactionTime     `fixed:"84,95,right"`
	Unknown3      string              `fixed:"96,96"`
	Unknown4      string              `fixed:"100,100"`
	DayOfEvent    string              `fixed:"103,110"`
	Unknown5      int                 `fixed:"123,123"`
	Splits        []*HY3Splits
	DQDescription *HY3DQDescription
}

type HY3TimeCode string

const (
	TimeCodeNormal     HY3TimeCode = " "
	TimeCodeScratch    HY3TimeCode = "S"
	TimeCodeNoShow     HY3TimeCode = "R"
	TimeCodeFalseStart HY3TimeCode = "F"
)

type HY3RelayEventEntryInfo struct {
	HY3Line     `fixed:"1,2"`
	TeamAbbr    string  `fixed:"3,7"`
	RelayTeam   string  `fixed:"8,8"`
	Gender      Gender  `fixed:"13,13"`
	Gender1     Gender  `fixed:"14,14"`
	Gender2     Gender  `fixed:"15,15"`
	Distance    int     `fixed:"18,21"`
	Stroke      string  `fixed:"22,22"`
	AgeLower    string  `fixed:"23,25"`
	AgeUpper    string  `fixed:"26,28"`
	EventFee    HY3Time `fixed:"33,38"`
	EventNumber string  `fixed:"39,41"`
	SeedTime1   string  `fixed:"44,44"`
	SeedCourse1 string  `fixed:"51,51"`
	SeedTime2   string  `fixed:"53,59"`
	SeedCourse2 string  `fixed:"60,60"`
}

type HY3RelayEventResults struct {
	HY3Line      `fixed:"1,2"`
	Type         string         `fixed:"3,3"`
	Time         string         `fixed:"4,11"`
	LengthUnit   string         `fixed:"12,12"`
	TimeCode     string         `fixed:"13,15"`
	Unknown1     string         `fixed:"16,20"`
	Heat         int            `fixed:"21,23"`
	Lane         int            `fixed:"24,26"`
	PlaceInHeat  int            `fixed:"27,29"`
	PlaceOVerall int            `fixed:"30,33"`
	Time1        HY3PlungerTime `fixed:"37,44"`
	Time2        HY3PlungerTime `fixed:"45,52"`
	Time3        HY3PlungerTime `fixed:"53,60"`
	Time4        HY3PlungerTime `fixed:"66,73"`
	Time5        HY3PlungerTime `fixed:"75,82"`
	DayOfEvent   string         `fixed:"103,110"`
}

type HY3DefaultTime float32

func (t HY3DefaultTime) MarshalTextFixedWidth() ([]byte, error) {
	if t == 0 {
		return []byte("0"), nil
	}
	return []byte(fmt.Sprintf("%.02f", t)), nil
}

func (t *HY3DefaultTime) UnmarshalTextFixedWidth(b []byte) error {
	if len(b) == 0 {
		*t = 0
		return nil
	}
	_, err := fmt.Sscan(string(b), t)
	return err
}

type HY3Time float32

func (t HY3Time) MarshalTextFixedWidth() ([]byte, error) {
	return []byte(fmt.Sprintf("%.02f", t)), nil
}

func (t *HY3Time) UnmarshalTextFixedWidth(b []byte) error {
	if len(b) == 0 {
		*t = 0
		return nil
	}
	_, err := fmt.Sscan(string(b), t)
	return err
}

func (t HY3Time) MarshalCSV() ([]byte, error) {
	return []byte(fmt.Sprintf("%v:%v.%v", int(t)/60, int(t)%60, t-HY3Time(int(t)))), nil
}

func (t *HY3Time) UnmarshalCSV(b []byte) error {
	if string(b) == "" {
		return nil
	}
	var minutes, seconds, hundreths int
	n, err := fmt.Sscanf(string(b), "%d:%d.%d", &minutes, &seconds, &hundreths)
	if err != nil {
		return fmt.Errorf("failed to parse string (%v); %v", string(b), err)
	}
	if n < 3 {
		return fmt.Errorf("wrong number of numbers parsed (%v): want(3), got(%v)", string(b), n)
	}
	*t = HY3Time(minutes*60+seconds) + HY3Time(hundreths)/100
	return nil
}

func (h HY3Time) String() string {
	if h == 0 {
		return "NT"
	}
	if int(h/60) > 0 {
		return fmt.Sprintf("%v:%05.2f", int(h/60), float32(h)-float32(int(h/60)*60))
	}
	return fmt.Sprintf("%04.2f", h)
}

type HY3PlungerTime struct {
	HY3Time
}

func (t HY3PlungerTime) MarshalTextFixedWidth() ([]byte, error) {
	if t.HY3Time == 0 {
		return nil, nil
	}
	return []byte(fmt.Sprintf("%.02f", t.HY3Time)), nil
}

type HY3ReactionTime float64

func (t HY3ReactionTime) MarshalTextFixedWidth() ([]byte, error) {
	if t == 0 {
		return nil, nil
	}
	if t > 10 {
		return []byte(fmt.Sprintf("0%v", int(t))), nil
	}
	return []byte(fmt.Sprintf("%.010f", t)), nil
}

func (t *HY3ReactionTime) UnmarshalTextFixedWidth(b []byte) error {
	if len(b) == 0 {
		*t = 0
		return nil
	}
	_, err := fmt.Sscan(string(b), t)
	return err
}

type HY3RelayEventLineUp struct {
	HY3Line  `fixed:"1,2"`
	Gender1  Gender `fixed:"3,3"`
	ID1      int    `fixed:"4,8"`
	Abbr1    string `fixed:"9,13"`
	Gender1X Gender `fixed:"14,14"`
	Leg1     int    `fixed:"15,15"`
	Gender2  Gender `fixed:"16,16"`
	ID2      int    `fixed:"17,21"`
	Abbr2    string `fixed:"22,26"`
	Gender2X Gender `fixed:"27,27"`
	Leg2     int    `fixed:"28,28"`
	Gender3  Gender `fixed:"29,29"`
	ID3      int    `fixed:"30,34"`
	Abbr3    string `fixed:"35,39"`
	Gender3X Gender `fixed:"40,40"`
	Leg3     int    `fixed:"41,41"`
	Gender4  Gender `fixed:"42,42"`
	ID4      int    `fixed:"43,47"`
	Abbr4    string `fixed:"48,52"`
	Gender4X Gender `fixed:"53,53"`
	Leg4     int    `fixed:"54,54"`
}

type HY3Splits struct {
	HY3Line `fixed:"1,2"`
	Times   HY3SplitTimes `fixed:"3,124"`
}

type HY3SplitTimes []*HY3SplitTime

func (h HY3SplitTimes) MarshalTextFixedWidth() ([]byte, error) {
	var buf bytes.Buffer
	enc := fixedwidth.NewEncoder(&buf)
	enc.SetLineTerminator([]byte{})
	for _, v := range h {
		if err := enc.Encode(v); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (h *HY3SplitTimes) UnmarshalTextFixedWidth(b []byte) error {
	for i := 0; i < len(b); i += 11 {
		split := &HY3SplitTime{}
		if err := fixedwidth.Unmarshal(b[i:i+11], split); err != nil {
			return err
		}
		*h = append(*h, split)
	}
	return nil
}

type HY3SplitTime struct {
	Type   string  `fixed:"1,1"`
	Length int     `fixed:"2,3,right"`
	Time   HY3Time `fixed:"4,11,right"`
}

type HY3DQDescription struct {
	HY3Line     `fixed:"1,2"`
	Code        string `fixed:"3,4"`
	Description string `fixed:"5,128"`
}

type checksumWriter struct {
	buf  *bytes.Buffer
	line []byte
}

func newChecksumWriter() *checksumWriter {
	return &checksumWriter{
		buf:  &bytes.Buffer{},
		line: make([]byte, 132),
	}
}

func (c *checksumWriter) reset() {
	for i := range c.line {
		c.line[i] = ' '
	}
	c.line[130] = '\r'
	c.line[131] = '\n'
}

func (c *checksumWriter) Write(b []byte) (int, error) {
	c.reset()
	copy(c.line, b)
	chk := hy3GenerateChecksum(c.line)
	c.line[128], c.line[129] = chk[0], chk[1]
	_, err := c.buf.Write(c.line)
	return len(b), err
}

func (c *checksumWriter) Read(p []byte) (int, error) {
	return c.buf.Read(p)
}

func hy3GenerateChecksum(s []byte) []byte {
	sum := 0
	for i, v := range s[:128] {
		if i%2 == 0 {
			sum += int(v)
		} else {
			sum += int(v) * 2
		}
	}
	sum /= 21
	sum += 205
	chk := fmt.Sprintf("%v%v", (sum % 10), (sum/10)%10)
	return []byte(chk)
}

func ParseHY3File(r io.Reader) (*HY3, error) {
	scanner := bufio.NewScanner(r)
	ret := &HY3{}
	var currTeam *HY3SwimTeam
	var currSwimmer *HY3Swimmer
	var currIndividualEntry *HY3IndividualEventEntryInfo
	var currResult *HY3IndividualEventResults
	for scanner.Scan() {
		line := scanner.Text()
		switch line[0:2] {
		case "A1":
			ret.FileDescriptor = &HY3FileDescriptor{}
			if err := fixedwidth.Unmarshal([]byte(line), ret.FileDescriptor); err != nil {
				return nil, err
			}
		case "B1":
			ret.MeetInfo = &HY3MeetInfo{}
			if err := fixedwidth.Unmarshal([]byte(line), ret.MeetInfo); err != nil {
				return nil, err
			}
		case "B2":
			ret.MeetAddress = &HY3MeetAddress{}
			if err := fixedwidth.Unmarshal([]byte(line), ret.MeetAddress); err != nil {
				return nil, err
			}
		case "B3":
			ret.MeetContact = &HY3MeetContact{}
			if err := fixedwidth.Unmarshal([]byte(line), ret.MeetContact); err != nil {
				return nil, err
			}
		case "C1":
			currTeam = &HY3SwimTeam{
				Name: &HY3SwimTeamNameInfo{},
			}
			ret.Teams = append(ret.Teams, currTeam)
			if err := fixedwidth.Unmarshal([]byte(line), currTeam.Name); err != nil {
				return nil, err
			}
		case "C2":
			if currTeam == nil {
				return nil, fmt.Errorf("TeamAddress before Team info")
			}
			currTeam.Address = &HY3SwimTeamAddressInfo{}
			if err := fixedwidth.Unmarshal([]byte(line), currTeam.Address); err != nil {
				return nil, err
			}
		case "C3":
			if currTeam == nil {
				return nil, fmt.Errorf("TeamContact before Team info")
			}
			currTeam.Contact = &HY3SwimTeamContactInfo{}
			if err := fixedwidth.Unmarshal([]byte(line), currTeam.Contact); err != nil {
				return nil, err
			}
		case "D1":
			currSwimmer = &HY3Swimmer{
				Info1: &HY3SwimmerInfo1{},
			}
			if currTeam == nil {
				return nil, fmt.Errorf("swimmerInfo before Team info")
			}
			currTeam.Swimmers = append(currTeam.Swimmers, currSwimmer)
			if err := fixedwidth.Unmarshal([]byte(line), currSwimmer.Info1); err != nil {
				return nil, err
			}
		case "D2":
		case "D3":
		case "D4":
		case "D5":
			continue
		case "E1":
			if currSwimmer == nil {
				return nil, fmt.Errorf("E1 before Swimmer info")
			}
			currIndividualEntry = &HY3IndividualEventEntryInfo{}
			if err := fixedwidth.Unmarshal([]byte(line), currIndividualEntry); err != nil {
				return nil, err
			}
			currSwimmer.IndividualEntries = append(currSwimmer.IndividualEntries, currIndividualEntry)
		case "E2":
			if currSwimmer == nil {
				return nil, fmt.Errorf("E2 before E1")
			}
			currResult = &HY3IndividualEventResults{}
			if err := fixedwidth.Unmarshal([]byte(line), currResult); err != nil {
				return nil, err
			}
			currIndividualEntry.Result = currResult
		case "G1":
			if currResult == nil {
				return nil, fmt.Errorf("G1 before Result info")
			}
			split := &HY3Splits{}
			if err := fixedwidth.Unmarshal([]byte(line), split); err != nil {
				return nil, err
			}
			currResult.Splits = append(currResult.Splits, split)
		case "H1":
			if currResult == nil {
				return nil, fmt.Errorf("H1 before Result info")
			}
			currResult.DQDescription = &HY3DQDescription{}
			if err := fixedwidth.Unmarshal([]byte(line), currResult.DQDescription); err != nil {
				return nil, err
			}
		}
	}
	return ret, nil
}

func GenerateHY3File(m *HY3, w io.Writer) error {
	buf := newChecksumWriter()
	enc := fixedwidth.NewEncoder(buf)
	enc.SetLineTerminator([]byte{})
	encode := func(v typeSetter, t string) error {
		v.setType(t)
		return enc.Encode(v)
	}
	if m.FileDescriptor != nil {
		if err := encode(m.FileDescriptor, "A1"); err != nil {
			return err
		}
	}
	if m.MeetInfo != nil {
		if err := encode(m.MeetInfo, "B1"); err != nil {
			return err
		}
	}
	if m.MeetAddress != nil {
		if err := encode(m.MeetAddress, "B2"); err != nil {
			return err
		}
	}
	if m.MeetContact != nil {
		if err := encode(m.MeetContact, "B3"); err != nil {
			return err
		}
	}
	for _, team := range m.Teams {
		if team.Name == nil {
			continue
		}
		if err := encode(team.Name, "C1"); err != nil {
			return err
		}
		if team.Address != nil {
			if err := encode(team.Address, "C2"); err != nil {
				return err
			}
		}
		if team.Contact != nil {
			if err := encode(team.Contact, "C3"); err != nil {
				return err
			}
		}
		for _, swimmer := range team.Swimmers {
			if err := encode(swimmer.Info1, "D1"); err != nil {
				return err
			}
			for _, entry := range swimmer.IndividualEntries {
				if err := encode(entry, "E1"); err != nil {
					return err
				}
				result := entry.Result
				if result == nil {
					continue
				}
				if err := encode(result, "E2"); err != nil {
					return err
				}
				for _, v := range result.Splits {
					if err := encode(v, "G1"); err != nil {
						return err
					}
				}
				if result.DQDescription != nil {
					if err := encode(result.DQDescription, "H1"); err != nil {
						return err
					}
				}
			}
		}
	}
	io.Copy(w, buf)
	return nil
}
