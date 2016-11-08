package hgtealib

import (
	"bytes"
	"errors"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Flush float64

const (
	First     Flush = 1.0
	InBetween Flush = 1.5
	Second    Flush = 2.0
	Monsoon   Flush = 3.0
	Autumn    Flush = 4.0
)

var flushes = []string{"First", "InBetween", "Second", "Monsoon", "Autumn"}

func (f Flush) String() string {
	if f == First {
		return "First"
	} else if f == InBetween {
		return "InBetween"
	} else if f == Second {
		return "Second"
	} else if f == Monsoon {
		return "Monsoon"
	} else if f == Autumn {
		return "Autumn"
	}
	return ""
}

// Timestamp       Date    Time    Tea     Rating  Comments        Pictures        Steep Time      Steeping Vessel Steep Temperature       Session Instance        Fixins
type Entry struct {
	Id                  int // TODO: Id => TeaId
	DateTime            time.Time
	Rating              int
	Comments            string
	SteepTime           time.Duration
	SteepingVessel      int
	SteepingTemperature int
	SessionInstance     string
	Fixins              []string
}

func (e *Entry) ParseDateTime(d, t string) error {
	if d == "" {
		return errors.New("Date is empty")
	}
	if t == "" {
		return errors.New("Time is empty")
	}

	// Determine the date
	darr := strings.Split(d, "/")

	month, err := strconv.Atoi(darr[0])
	if err != nil {
		log.Println(err)
	}

	day, err := strconv.Atoi(darr[1])
	if err != nil {
		log.Println(err)
	}

	year, err := strconv.Atoi(darr[2])
	if err != nil {
		log.Println(err)
	}

	// Determine the time
	minute, err := strconv.Atoi(t[len(t)-2:])
	if err != nil {
		log.Println(err)
	}

	hour, err := strconv.Atoi(t[:len(t)-2])
	if err != nil {
		log.Println(err)
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println(err)
	}

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

func (t *Tea) Add(entry Entry) error {
	t.log[entry.DateTime] = entry
	t.logSortedKeys = append(t.logSortedKeys, entry.DateTime)
	sort.Sort(t.logSortedKeys)

	return nil
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

		t.median = ratings[((len(t.log) + 1) / 2)]
		// } else if t.median == 0 && len(t.log) == 1 {
		// t.median = t.log[
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

func (t *Tea) String() string {
	var buf bytes.Buffer

	if t.Picked.Year > 0 {
		buf.WriteString(strconv.Itoa(t.Picked.Year))
		buf.WriteString(" ")
	}
	buf.WriteString(t.Name)

	return buf.String()
}
