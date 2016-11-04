package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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

	// log.Printf("Found %d teas\n", len(db))
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

	// log.Printf("Found %d journal entries\n", len(journal))
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

func main() {
	statsCmd := flag.NewFlagSet("stats", flag.ExitOnError)

	if os.Args[1] == "stats" {
		stockedFlag := statsCmd.Bool("stocked", false, "Only display stocked teas")
		samplesFlag := statsCmd.Bool("samples", false, "Only display tea samples")
		teaTypes := statsCmd.String("types", "", "Comma-delimited list of tea types to select")

		statsCmd.Parse(os.Args[2:])

		typesFilter := make(map[string]struct{})
		if len(*teaTypes) > 0 {
			for _, typeFilter := range strings.Split(*teaTypes, ",") {
				typesFilter[strings.ToLower(typeFilter)] = struct{}{}
			}
		}

		db, err := buildDatabase(*stockedFlag, *samplesFlag, typesFilter)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%-60s %6s %6s %6s %6s\n", "Name", "Num", "Avg", "Median", "Mode")
		for _, tea := range db {
			fmt.Printf("%-60s %6d %6d %6d %6d\n", tea.String(), len(tea.Log), tea.Average(), tea.Median(), tea.Mode())
		}
	}
}
