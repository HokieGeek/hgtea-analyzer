package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

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
		// } else if os.Args[1] == "journal" || os.Args[1] == "log" {
		// db, err := buildDatabase(*stockedFlag, *samplesFlag, typesFilter)
		// if err != nil {
		// log.Fatal(err)
		// }

		// if len(os.Args) > 2 {
		// for _,entries := range db[os.Args[2] {
		// fmt.Println(entries)
		// }
		// }
	}
}
