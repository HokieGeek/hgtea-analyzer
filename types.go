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

// Timestamp       Date    Time    Tea     Rating  Comments        Pictures        Steep Time      Steeping Vessel Steep Temperature       Session Instance        Fixins
type Entry struct {
	Id                  int
	Date                string // TODO
	Time                string // TODO
	DateTime            time.Time
	Rating              int
	Comments            string
	SteepTime           time.Duration
	SteepingVessel      int
	SteepingTemperature int
	SessionInstance     string
	Fixins              []string
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

type TeaPickPeriod struct {
	Year  int
	Flush float64
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
	log.Printf("%d\n", len(t.logSortedKeys))
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

func getEntryTime(d string, t string) (time.Time, error) {
	if d == "" {
		return time.Now(), errors.New("Date is empty")
	}
	if t == "" {
		return time.Now(), errors.New("Time is empty")
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

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, loc), nil
}

func getEntryDuration(d string) time.Duration {
	re := regexp.MustCompile("[[:space:]]")
	dur, err := time.ParseDuration(re.ReplaceAllString(d, ""))
	if err != nil {
		return time.Nanosecond
	}
	return dur
}

func newEntry(entry []string) (*Entry, error) {
	if len(entry) < 11 {
		return nil, errors.New("Invalid data")
	}

	e := new(Entry)

	// e.Date = entry[1]
	// e.Time = entry[2]

	e.Id, _ = strconv.Atoi(entry[3])
	dateTime, err := getEntryTime(entry[1], entry[2])
	if err != nil {
		return nil, err
	}
	e.DateTime = dateTime

	e.Rating, _ = strconv.Atoi(entry[4])
	e.Comments = entry[5]

	e.SteepTime = getEntryDuration(entry[7])

	e.SteepingVessel, _ = strconv.Atoi(entry[8])
	e.SteepingTemperature, _ = strconv.Atoi(entry[9])
	e.SessionInstance = entry[10]
	e.Fixins = strings.Split(entry[11], ";")

	return e, nil
}

func newTea(data []string) (*Tea, error) {
	if len(data) < 22 {
		return nil, errors.New("Data badly formatted")
	}

	t := new(Tea)

	var err error
	t.Id, err = strconv.Atoi(data[2])
	if err != nil {
		return nil, err
	}
	t.Name = data[3]
	t.Type = data[4]
	t.Size = data[18]
	t.LeafGrade = data[15]

	t.Origin.Country = data[14]
	t.Origin.Region = data[5]

	t.Storage.Stocked = (data[19] == "TRUE")
	t.Storage.Aging = (data[20] == "TRUE")

	if data[6] != "" {
		t.Picked.Year, err = strconv.Atoi(data[6])
		if err != nil {
			return nil, err
		}
	}
	if data[7] != "" {
		t.Picked.Flush, err = strconv.ParseFloat(data[7], 64)
		if err != nil {
			return nil, err
		}
	}

	t.Purchased.Location = data[8]
	t.Purchased.Date = data[9]
	if data[10] != "" {
		t.Purchased.Price, err = strconv.ParseFloat(data[10], 64)
		if err != nil {
			return nil, err
		}
	}
	if data[21] != "" {
		t.Purchased.Packaging, err = strconv.Atoi(data[21])
		if err != nil {
			return nil, err
		}
	}

	t.log = make(map[time.Time]Entry)
	t.logSortedKeys = make(TimeSlice, 0)

	return t, nil
}
