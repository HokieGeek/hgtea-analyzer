package main

import (
	"bytes"
	"errors"
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
}

func (t *Tea) addEntry(entry []string) error {
	return nil
}

func (t *Tea) String() string {
	var buf bytes.Buffer

	buf.WriteString("[")
	buf.WriteString(strconv.Itoa(t.Id))
	buf.WriteString("] ")
	buf.WriteString(t.Name)

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

	return t, nil
}
