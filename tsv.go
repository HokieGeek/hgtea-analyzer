package hgtealib

import (
	"encoding/csv"
	"errors"
	"golang.org/x/net/proxy"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getSheetTsv(url, proxyAddr string) ([][]string, error) {
	var response *http.Response
	var err error
	if proxyAddr != "" {
		// TODO: Determine if html or socks5 based on the protocol: http:// , socks5://
		socks5Proxy := proxyAddr
		dialer, err := proxy.SOCKS5("tcp", socks5Proxy, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}
		httpTransport.Dial = dialer.Dial

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		response, err = httpClient.Do(req)
	} else {
		response, err = http.Get(url)
	}
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	r := csv.NewReader(response.Body)
	r.Comma = '\t'
	db, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newEntryFromTsv(entry []string) (*Entry, error) {
	if len(entry) < 11 {
		return nil, errors.New("Invalid data")
	}

	e := new(Entry)

	e.Tea, _ = strconv.Atoi(entry[3])
	e.ParseDateTime(entry[1], entry[2])

	e.Rating, _ = strconv.Atoi(entry[4])
	e.Comments = entry[5]

	e.ParseSteepTime(entry[7])
	dummy_int, _ := strconv.Atoi(entry[8])
	e.SteepingVessel = VesselType(dummy_int)
	// e.SteepingVessel, _ = strconv.Atoi(entry[8])
	e.SteepingTemperature, _ = strconv.Atoi(entry[9])
	if e.SteepingTemperature == 0 {
		// TODO: make this value depend on the type (if green or oolong, for example)
		e.SteepingTemperature = 212
	}

	e.SessionInstance = entry[10]
	for _, f := range strings.Split(entry[11], ";") {
		if f != "" {
			dummy, _ := strconv.Atoi(f)
			e.Fixins = append(e.Fixins, TeaFixin(dummy))
		}
	}

	return e, nil
}

func newTeaFromTsv(data []string) (*Tea, error) {
	if len(data) < 22 {
		return nil, errors.New("Data badly formatted")
	}

	t := new(Tea)
	t.log = make(map[time.Time]Entry)
	t.logSortedKeys = make(TimeSlice, 0)

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
		dummy_float, err := strconv.ParseFloat(data[7], 64)
		if err != nil {
			return nil, err
		}
		t.Picked.Flush = Flush(dummy_float)
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

	return t, nil
}

func NewFromTsv(teas_url, log_url, proxyAddr string) (*HgTeaDb, error) {
	// Get the tea database
	teasTsv, err := getSheetTsv(teas_url, proxyAddr)
	if err != nil {
		return nil, err
	}

	teas := make([]*Tea, 0)
	for _, tea := range teasTsv[1:] {
		t, err := newTeaFromTsv(tea)
		if err != nil {
			return nil, err
		}

		teas = append(teas, t)
	}

	// Add the journal entries
	journal, err := getSheetTsv(log_url, proxyAddr)
	if err != nil {
		return nil, err
	}

	entries := make([]*Entry, 0)
	for _, entry := range journal[1:] {
		e, err := newEntryFromTsv(entry)
		if err != nil {
			return nil, err
		}

		entries = append(entries, e)
	}

	return newHgTeaDb(teas, entries)
}
