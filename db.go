package hgtealib

import (
	"encoding/csv"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type HgTeaDb struct {
	teas          map[int]Tea
	log           map[time.Time]Entry
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

func (d *HgTeaDb) Log(filter *Filter) ([]Entry, error) {
	log := make([]Entry, 0)
	for _, k := range d.logSortedKeys {
		log = append(log, d.log[k])
	}
	return log, nil
}

/*
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

	// db.log = make([]*Entry, 0)
	db.log = make(map[time.Time]Entry)
	for _, entry := range journal[1:] {
		e, err := newEntry(entry)
		if err != nil {
			return nil, err
		}

		db.log[e.DateTime] = *e
		db.logSortedKeys = append(db.logSortedKeys, e.DateTime)
		sort.Sort(db.logSortedKeys)

		id, _ := strconv.Atoi(entry[3])
		if tea, ok := db.teas[id]; ok {
			err := tea.Add(*e)
			db.teas[id] = tea // TODO: why do I have to do this?
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
