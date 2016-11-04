package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
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

func buildDatabase() ([]Tea, error) {
	// db_url := "https://docs.google.com/spreadsheets/d/1-U45bMxRE4_n3hKRkTPTWHTkVKC8O3zcSmkjEyYFYOo/pub?output=tsv"
	// journal_url := "https://docs.google.com/spreadsheets/d/1pHXWycR9_luPdHm32Fb2P1Pp7l29Vni3uFH_q3TsdbU/pub?output=tsv"

	db, err := getSheet("https://docs.google.com/spreadsheets/d/1-U45bMxRE4_n3hKRkTPTWHTkVKC8O3zcSmkjEyYFYOo/pub?output=tsv")
	if err != nil {
		return nil, err
	}

	teas := make([]Tea, len(db))

	log.Printf("Found %d records\n", len(db))
	for _, entry := range db {
		t, _ := newTea(entry)
		fmt.Println(t)
	}

	return teas, nil
}

func main() {
	buildDatabase()
}
