package hgtealib

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"
)

type HgTeaDb struct {
	teas          map[int]Tea
	log           []*Entry
	logSortedKeys TimeSlice
}

type Filter struct {
	stockedOnly bool
	samplesOnly bool
	types       map[string]struct{}
}

func (f *Filter) StockedOnly() *Filter {
	f.stockedOnly = true
	return f
}

func (f *Filter) SamplesOnly() *Filter {
	f.samplesOnly = true
	return f
}

func (f *Filter) Types(v []string) *Filter {
	for _, t := range v {
		f.Type(t)
	}
	return f
}

func (f *Filter) Type(v string) *Filter {
	f.types[strings.ToLower(v)] = struct{}{}
	return f
}

func NewFilter() *Filter {
	f := new(Filter)

	f.stockedOnly = false
	f.samplesOnly = false
	f.types = make(map[string]struct{})

	return f
}

func (d *HgTeaDb) Teas(filter *Filter) (map[int]Tea, error) {
	teas := make(map[int]Tea)
	for k, v := range d.teas {
		// Now apply the filters
		if filter.stockedOnly && !v.Storage.Stocked {
			continue
		}

		// if filter.samplesOnly && !strings.Contains(strings.ToLower(v.Size), "sample") {
		// continue
		// }

		if len(filter.types) > 0 {
			if _, ok := filter.types[strings.ToLower(v.Type)]; !ok {
				continue
			}
		}

		teas[k] = v
	}
	return teas, nil
}

/*
func (d *HgTeaDb) Log(filter *Filter) (map[time.Time]*Entry, error) {
	teas := make(map[time.Time]Tea)
	for k, v := range d.log {
		log[k] = v
	}
	return log, nil
}
*/

func New(teas_url string, log_url string) (*HgTeaDb, error) {
	db := new(HgTeaDb)

	// Get the tea database
	teas, err := getSheet(teas_url)
	if err != nil {
		return nil, err
	}

	db.teas = make(map[int]Tea)
	for _, tea := range teas[1:] {
		t, err := newTea(tea)
		if err != nil {
			return nil, err
		}

		db.teas[t.Id] = *t
	}

	// Add the journal entries
	journal, err := getSheet(log_url)
	if err != nil {
		return nil, err
	}

	db.log = make([]*Entry, 0)
	for _, entry := range journal[1:] {
		id, _ := strconv.Atoi(entry[3])
		// TODO: also, add it to db.log
		if tea, ok := db.teas[id]; ok {
			err := tea.Add(entry)
			if err != nil {
				return nil, err
			}
		}
	}

	return db, nil
}

func getSheet(url string) ([][]string, error) {
	response, err := http.Get(url)
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
