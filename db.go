package main

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"
)

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

func buildDatabase(stockedOnly bool, samplesOnly bool, types map[string]struct{}) (map[int]Tea, error) {
	// Get the tea database
	db, err := getSheet("https://docs.google.com/spreadsheets/d/1-U45bMxRE4_n3hKRkTPTWHTkVKC8O3zcSmkjEyYFYOo/pub?output=tsv")
	if err != nil {
		return nil, err
	}

	teas := make(map[int]Tea)
	for _, tea := range db[1:] {
		t, err := newTea(tea)
		if err != nil {
			return nil, err
		}

		// Now apply the filters
		if stockedOnly && !t.Stocked {
			continue
		}

		if samplesOnly && !strings.Contains(strings.ToLower(t.Size), "sample") {
			continue
		}

		if len(types) > 0 {
			if _, ok := types[strings.ToLower(t.Type)]; !ok {
				continue
			}
		}

		teas[t.Id] = *t
	}

	// Add the journal entries
	journal, err := getSheet("https://docs.google.com/spreadsheets/d/1pHXWycR9_luPdHm32Fb2P1Pp7l29Vni3uFH_q3TsdbU/pub?output=tsv")
	if err != nil {
		return nil, err
	}

	for _, entry := range journal[1:] {
		id, _ := strconv.Atoi(entry[3])
		if tea, ok := teas[id]; ok {
			err := tea.Add(entry)
			if err != nil {
				return nil, err
			}
		}
	}

	return teas, nil
}
