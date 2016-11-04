package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func buildDatabase() (map[int]Tea, error) {
	// Get the tea database
	db, err := getSheet("https://docs.google.com/spreadsheets/d/1-U45bMxRE4_n3hKRkTPTWHTkVKC8O3zcSmkjEyYFYOo/pub?output=tsv")
	if err != nil {
		return nil, err
	}

	log.Printf("Found %d teas\n", len(db))
	teas := make(map[int]Tea)
	for _, tea := range db[1:] {
		t, err := newTea(tea)
		if err != nil {
			return nil, err
		}
		teas[t.Id] = *t
	}

	// Add the journal entries
	journal, err := getSheet("https://docs.google.com/spreadsheets/d/1pHXWycR9_luPdHm32Fb2P1Pp7l29Vni3uFH_q3TsdbU/pub?output=tsv")
	if err != nil {
		return nil, err
	}

	log.Printf("Found %d journal entries\n", len(journal))
	for _, entry := range journal[1:] {
		id, _ := strconv.Atoi(entry[3])
		tea := teas[id]
		err := tea.Add(entry)
		if err != nil {
			return nil, err
		}
	}

	return teas, nil
}

func main() {
	db, err := buildDatabase()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-50s %6s %6s %6s\n", "Name", "Avg", "Median", "Mode")
	for _, tea := range db {
		fmt.Printf("%-50s %6d %6d %6d\n", tea.Name, tea.Average(), tea.Median(), tea.Mode())
	}
}
