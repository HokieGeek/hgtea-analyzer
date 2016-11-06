package main

import (
	"flag"
	"fmt"
	"github.com/hokiegeek/hgtealib"
	"log"
	// "os"
	"strings"
)

func main() {
	// commonOpts := flag.NewFlagSet("common", flag.ExitOnError)

	stockedFlag := flag.Bool("stocked", false, "Only display stocked teas")
	samplesFlag := flag.Bool("samples", false, "Only display tea samples")
	teaTypes := flag.String("types", "", "Comma-delimited list of tea types to select")

	flag.Parse()

	typesFilter := make(map[string]struct{})
	if len(*teaTypes) > 0 {
		for _, typeFilter := range strings.Split(*teaTypes, ",") {
			typesFilter[strings.ToLower(typeFilter)] = struct{}{}
		}
	}

	db, err := hgtealib.BuildDatabase(*stockedFlag, *samplesFlag, typesFilter)
	if err != nil {
		log.Fatal(err)
	}

	command := flag.Arg(0)

	// if os.Args[1] == "stats" {
	if command == "stats" {
		// statsCmd := flag.NewFlagSet("stats", flag.ExitOnError)

		// statsCmd.Parse(os.Args[2:])
		fmt.Printf("%-60s %6s %6s %6s %6s\n", "Name", "Num", "Avg", "Median", "Mode")
		for _, tea := range db {
			fmt.Printf("%-60s %6d %6d %6d %6d\n", tea.String(), len(tea.Log), tea.Average(), tea.Median(), tea.Mode())
		}
		// } else if os.Args[1] == "journal" || os.Args[1] == "log" {
	} else if command == "teas" {
		fmt.Printf("%-60s %s\n", "Name", "Year")
		for _, tea := range db {
			fmt.Printf("%-60s %d\n", tea.Name, tea.Picked.Year)
		}
	} else if command == "log" {
		fmt.Println("TODO")
	}
}
