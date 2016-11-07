package main

import (
	"flag"
	"fmt"
	"github.com/hokiegeek/hgtealib"
	"log"
	"regexp"
	"strings"
	"time"
)

func main() {
	teas_url := "https://docs.google.com/spreadsheets/d/1-U45bMxRE4_n3hKRkTPTWHTkVKC8O3zcSmkjEyYFYOo/pub?output=tsv"
	log_url := "https://docs.google.com/spreadsheets/d/1pHXWycR9_luPdHm32Fb2P1Pp7l29Vni3uFH_q3TsdbU/pub?output=tsv"

	stockedFlag := flag.Bool("stocked", false, "Only display stocked teas")
	// samplesFlag := flag.Bool("samples", false, "Only display tea samples")
	teaTypes := flag.String("types", "", "Comma-delimited list of tea types to select")

	flag.Parse()

	filter := hgtealib.NewFilter()
	if *stockedFlag {
		filter.StockedOnly()
	}
	// if *samplesFlag {
	// filter.SamplesOnly()
	// }
	filter.Types(strings.Split(*teaTypes, ","))

	db, err := hgtealib.New(teas_url, log_url)
	if err != nil {
		log.Fatal(err)
	}

	command := flag.Arg(0)

	if command == "teas" {
		fmt.Printf("%-60s %4s %6s %8s\n", "Name", "Year", "Entries", "Flush")
		teas, _ := db.Teas(filter)
		for _, tea := range teas {
			fmt.Printf("%-60s %d %7d %8f\n", tea.Name, tea.Picked.Year, tea.LogLen(), tea.Picked.Flush)
		}
	} else if command == "log" {
		// if len(flag.Args()) > 0 {
		// fmt.Printf("%v\n", flag.Args()[1:])
		logCmd := flag.NewFlagSet("log", flag.ExitOnError)
		noPrettyPrintFlag := logCmd.Bool("nopretty", false, "Formats the table prettily")

		logCmd.Parse(flag.Args()[1:])

		headerFmt := "%-21s\t%-60s\t%10s\t%s\n"
		entryFmt := "%s\t%-60s\t%10s\t%d\n"
		if *noPrettyPrintFlag {
			re := regexp.MustCompile("%-?[0-9]+")
			headerFmt = re.ReplaceAllString(headerFmt, "%")
			entryFmt = re.ReplaceAllString(entryFmt, "%")
		}

		/*
			// DateTime            time.Time
			// Rating              int
			Comments            string
			// SteepTime           time.Duration
			SteepingVessel      int
			SteepingTemperature int
			SessionInstance     string
			Fixins              []string
		*/
		fmt.Printf(headerFmt, "Time", "Tea", "Steep Time", "Rating")

		teas, _ := db.Teas(hgtealib.NewFilter())
		log, _ := db.Log(hgtealib.NewFilter())
		for _, v := range log {
			tea := teas[v.Id]
			fmt.Printf(entryFmt, v.DateTime.Format(time.RFC822Z), tea.String(), v.SteepTime, v.Rating)
		}
		/*
			tea := teas[0]
			for _, v := range tea.Log() {
				fmt.Printf(entryFmt, v.DateTime.Format(time.RFC822Z), tea.String(), v.SteepTime, v.Rating)
			}
		*/
		// entries, _ := db.Log(new(Filter))
		// for _, entry := range entries {
		// fmt.Println(entry)
		// fmt.Printf("%-60s %d %7d %8s %8s %17s\n", tea.Name, tea.Picked.Year, len(tea.Log), tea.Date, tea.Time, tea.DateTime)
		// }
	} else if command == "stats" {
		// statsCmd := flag.NewFlagSet("stats", flag.ExitOnError)

		fmt.Printf("%-60s %6s %6s %6s %6s\n", "Name", "Entries", "Avg", "Median", "Mode")
		teas, _ := db.Teas(filter)
		for _, tea := range teas {
			fmt.Printf("%-60s %6d %6d %6d %6d\n", tea.String(), tea.LogLen(), tea.Average(), tea.Median(), tea.Mode())
		}
	}
}
