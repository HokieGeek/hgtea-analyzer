package hgtealib

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
var TeaFlushTypes = [ ["Spring", "Summer", "Fall", "Winter"],
                      ["1st Flush", "2nd Flush", "Monsoon Flush", "Autumn Flush"] ];
var TeaFlushTypes_Std = 0;
var TeaFlushTypes_Indian = 1;
*/
type Flush float64

const (
	First     Flush = 1.0
	InBetween Flush = 1.5
	Second    Flush = 2.0
	Monsoon   Flush = 3.0
	Autumn    Flush = 4.0
)

func (f Flush) String() string {
	switch {
	case f == First:
		return "First"
	case f == InBetween:
		return "InBetween"
	case f == Second:
		return "Second"
	case f == Monsoon:
		return "Monsoon"
	case f == Autumn:
		return "Autumn"
	default:
		return ""
	}
}

type VesselType int

const (
	FrenchPress VesselType = 0 + iota
	ShipiaoYixing
	TeazerTumbler
	TeaStick
	MeshSpoon
	SaucePan
	Cup
	Bowl
	Gaiwan
	Other
)

func (v VesselType) String() string {
	switch {
	case v == FrenchPress:
		return "French Press"
	case v == ShipiaoYixing:
		return "Shipiao Yixing"
	case v == TeazerTumbler:
		return "Tea-zer Tumbler"
	case v == TeaStick:
		return "Tea stick"
	case v == MeshSpoon:
		return "Mesh spoon"
	case v == SaucePan:
		return "Sauce pan"
	case v == Cup:
		return "Cup"
	case v == Bowl:
		return "Bowl"
	case v == Gaiwan:
		return "Gaiwan"
	case v == Other:
		return "Other"
	default:
		return ""
	}
}

type TeaFixin int

const (
	Milk TeaFixin = 0 + iota
	Cream
	HalfAndHalf
	Sugar
	BrownSugar
	RawSugar
	Honey
	VanillaExtract
	VanillaBean
)

func (f TeaFixin) String() string {
	switch {
	case f == Milk:
		return "Milk"
	case f == Cream:
		return "Cream"
	case f == HalfAndHalf:
		return "Half & half"
	case f == Sugar:
		return "Sugar"
	case f == BrownSugar:
		return "Brown sugar"
	case f == RawSugar:
		return "Raw sugar"
	case f == Honey:
		return "Honey"
	case f == VanillaExtract:
		return "Vanilla extract"
	case f == VanillaBean:
		return "Vanilla bean"
	default:
		return ""
	}
}

/*
var TeaProductRatings = ["Value", "Leaf Aroma", "Brewed Aroma"];
var TeaPackagingTypes = ["Loose Leaf", "Bagged", "Tuo", "Beeng", "Brick", "Mushroom", "Square"];
*/

// Timestamp       Date    Time    Tea     Rating  Comments        Pictures        Steep Time      Steeping Vessel Steep Temperature       Session Instance        Fixins
type Entry struct {
	Tea                 int
	DateTime            time.Time
	Rating              int
	Comments            string
	SteepTime           time.Duration
	SteepingVessel      VesselType
	SteepingTemperature int
	SessionInstance     string
	Fixins              []TeaFixin
}

func (e *Entry) ParseDateTime(d, t string) error {
	// Validate the date field
	if d == "" {
		return errors.New("Date is empty")
	}

	darr := strings.Split(d, "/")
	if len(darr) < 3 || darr[0] == "" || darr[1] == "" || darr[2] == "" {
		return errors.New(fmt.Sprintf("Date field is invalid: %s", d))
	}

	// Validate the time field
	if t == "" {
		return errors.New("Time is empty")
	}

	if len(t) < 3 {
		return errors.New(fmt.Sprintf("Time field is invalid: %s", t))
	}

	// Determine the date
	month, err := strconv.Atoi(darr[0])
	if err != nil {
		return err
	}

	day, err := strconv.Atoi(darr[1])
	if err != nil {
		return err
	}

	year, err := strconv.Atoi(darr[2])
	if err != nil {
		return err
	}

	// Determine the time
	minute, err := strconv.Atoi(t[len(t)-2:])
	if err != nil {
		return err
	}

	hour, err := strconv.Atoi(t[:len(t)-2])
	if err != nil {
		return err
	}

	// In theory, this should never result in an error
	loc, _ := time.LoadLocation("America/New_York")

	e.DateTime = time.Date(year, time.Month(month), day, hour, minute, 0, 0, loc)

	return nil
}

func (e *Entry) ParseSteepTime(d string) error {
	re := regexp.MustCompile("[[:space:]]")
	dur, err := time.ParseDuration(re.ReplaceAllString(d, ""))
	if err != nil {
		return err
	}
	e.SteepTime = dur

	return nil
}

func (e *Entry) Equal(other *Entry) bool {
	return e.Tea == other.Tea &&
		e.DateTime.Equal(other.DateTime) &&
		e.Rating == other.Rating &&
		e.Comments == other.Comments &&
		e.SteepTime.Nanoseconds() == other.SteepTime.Nanoseconds() &&
		e.SteepingVessel == other.SteepingVessel &&
		e.SteepingTemperature == other.SteepingTemperature &&
		e.SessionInstance == other.SessionInstance &&
		(len(e.Fixins) == len(other.Fixins))
}

type TimeSlice []time.Time

func (e TimeSlice) Len() int {
	return len(e)
}

func (e TimeSlice) Less(i, j int) bool {
	return e[i].Before(e[j])
}

func (e TimeSlice) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Timestamp       Date    ID      Name    Type    Region  Year    Flush   Purchase Location       Purchase Date   Purchase Price  Ratings Comments        Pictures        Country Leaf Grade      Blended Teas   Blend Ratio     Size    Stocked Aging   Packaging
type TeaOrigin struct {
	Country string
	Region  string
}

func (o TeaOrigin) String() string {
	var buf bytes.Buffer

	if o.Region != "" {
		buf.WriteString(o.Region)
		buf.WriteString(", ")
	}
	buf.WriteString(o.Country)

	return buf.String()
}

type TeaPickPeriod struct {
	Year  int
	Flush Flush
}

type TeaStorageState struct {
	Stocked bool
	Aging   bool
}

type TeaPurchaseInfo struct {
	Location  string
	Date      string
	Price     float64
	Packaging int
}

type Tea struct {
	Id            int
	Name          string
	Type          string
	Picked        TeaPickPeriod
	Origin        TeaOrigin
	Storage       TeaStorageState
	Purchased     TeaPurchaseInfo
	Size          string
	LeafGrade     string // TODO: enum
	log           map[time.Time]Entry
	logSortedKeys TimeSlice
	average       int
	median        int
	mode          int
}

func (t *Tea) Add(entry Entry) {
	if t.log == nil {
		t.log = make(map[time.Time]Entry)
		t.logSortedKeys = make(TimeSlice, 0)
	}

	t.log[entry.DateTime] = entry
	t.logSortedKeys = append(t.logSortedKeys, entry.DateTime)
	sort.Sort(t.logSortedKeys)
}

func (t *Tea) Log() []Entry {
	log := make([]Entry, 0)
	for _, k := range t.logSortedKeys {
		log = append(log, t.log[k])
	}
	return log
}

func (t *Tea) LogLen() int {
	return len(t.log)
}

func (t *Tea) Average() int {
	if t.average == 0 && len(t.log) > 0 {
		var total int
		for _, entry := range t.log {
			total += entry.Rating
		}
		t.average = total / len(t.log)
	}

	return t.average
}

func (t *Tea) Median() int {
	if t.median == 0 && len(t.log) > 1 {
		ratings := make([]int, len(t.log))
		var count int
		for _, entry := range t.log {
			ratings[count] = entry.Rating
			count++
		}
		sort.Ints(ratings)

		if (len(ratings) % 2) == 0 {
			t.median = (ratings[len(ratings)/2] + ratings[(len(ratings)/2)-1]) / 2
		} else {
			t.median = ratings[len(ratings)/2]
		}
	}

	return t.median
}

func (t *Tea) Mode() int {
	if t.mode == 0 && len(t.log) > 0 {
		ratings := make([]int, 5)
		for _, entry := range t.log {
			ratings[entry.Rating]++
		}

		var max int
		for rating, count := range ratings {
			if count > ratings[max] {
				max = rating
			}
		}

		t.mode = max
	}

	return t.mode
}

func (t *Tea) Equal(other *Tea) bool {
	return t.Id == other.Id &&
		t.Name == other.Name &&
		t.Type == other.Type &&
		t.Picked.Year == other.Picked.Year &&
		t.Picked.Flush == other.Picked.Flush &&
		t.Size == other.Size &&
		t.Origin.Country == other.Origin.Country &&
		t.Origin.Region == other.Origin.Region &&
		t.Storage.Stocked == other.Storage.Stocked &&
		t.Storage.Aging == other.Storage.Aging &&
		t.Purchased.Location == other.Purchased.Location &&
		t.Purchased.Date == other.Purchased.Date &&
		t.Purchased.Price == other.Purchased.Price &&
		t.Purchased.Packaging == other.Purchased.Packaging &&
		t.LeafGrade == other.LeafGrade
	/*
		t.LeafGrade     string // TODO: enum

		t.log           map[time.Time]Entry
		t.logSortedKeys TimeSlice
		t.average       int
		t.median        int
		t.mode          int
	*/
}

func (t *Tea) String() string {
	var buf bytes.Buffer

	if t.Picked.Year > 0 {
		buf.WriteString(strconv.Itoa(t.Picked.Year))
		buf.WriteString(" ")
	}
	buf.WriteString(t.Name)

	return buf.String()
}
