package main

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Timestamp       Date    Time    Tea     Rating  Comments        Pictures        Steep Time      Steeping Vessel Steep Temperature       Session Instance        Fixins
type Entry struct {
	Date                string // TODO
	Time                string // TODO
	Rating              int
	Comments            string
	SteepTime           string
	SteepingVessel      int
	SteepingTemperature int
	SessionInstance     string
	Fixins              []string
}

func newEntry(entry []string) (*Entry, error) {
	if len(entry) < 11 {
		return nil, errors.New("Invalid data")
	}

	e := new(Entry)

	e.Date = entry[1]
	e.Time = entry[2]
	e.Rating, _ = strconv.Atoi(entry[4])
	e.Comments = entry[5]
	e.SteepTime = entry[7]
	e.SteepingVessel, _ = strconv.Atoi(entry[8])
	e.SteepingTemperature, _ = strconv.Atoi(entry[9])
	e.SessionInstance = entry[10]
	e.Fixins = strings.Split(entry[11], ";")

	return e, nil
}

// Timestamp       Date    ID      Name    Type    Region  Year    Flush   Purchase Location       Purchase Date   Purchase Price  Ratings Comments        Pictures        Country Leaf Grade      Blended Teas   Blend Ratio     Size    Stocked Aging   Packaging
type Tea struct {
	Id      int
	Name    string
	Type    string
	Stocked bool
	Aging   bool
	Size    string
	Log     map[string]Entry
	average int
	median  int
	mode    int
}

func (t *Tea) Add(entry []string) error {
	e, err := newEntry(entry)
	if err != nil {
		return err
	}

	ts := fmt.Sprintf("%s_%s", e.Date, e.Time)
	t.Log[ts] = *e

	return nil
}

func (t *Tea) Average() int {
	if t.average == 0 && len(t.Log) > 0 {
		var total int
		for _, entry := range t.Log {
			total += entry.Rating
		}
		t.average = total / len(t.Log)
	}

	return t.average
}

func (t *Tea) Median() int {
	if t.median == 0 && len(t.Log) > 0 {
		ratings := make([]int, len(t.Log))
		for _, entry := range t.Log {
			ratings = append(ratings, entry.Rating)
		}

		sort.Ints(ratings)
		t.median = ratings[((len(t.Log) + 1) / 2)]
	}

	return t.median
}

func (t *Tea) Mode() int {
	if t.mode == 0 && len(t.Log) > 0 {
		ratings := make(map[int]int)
		for _, entry := range t.Log {
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

	buf.WriteString("[")
	buf.WriteString(strconv.Itoa(t.Id))
	buf.WriteString("] ")
	buf.WriteString(t.Name)
	buf.WriteString(", ")
	buf.WriteString(strconv.Itoa(len(t.Log)))
	buf.WriteString(" entries")

	return buf.String()
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
	t.Stocked = (data[20] == "TRUE")
	t.Aging = (data[21] == "TRUE")
	t.Size = data[19]
	t.Log = make(map[string]Entry)

	t.average = 0
	t.median = 0
	t.mode = 0

	return t, nil
}