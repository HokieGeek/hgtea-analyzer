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

	proxyStr := flag.String("proxy", "", "Use the given proxy")

	stockedFlag := flag.Bool("stocked", false, "Only display stocked teas")
	// samplesFlag := flag.Bool("samples", false, "Only display tea samples")
	teaTypes := flag.String("types", "", "Comma-delimited list of tea types to select")

	// TODO: columnsStr := flag.String("columns", "", "Comma-delimited list of the columns to display")
	noPrettyPrintFlag := flag.Bool("nopretty", false, "Formats the table prettily")

	flag.Parse()

	filter := hgtealib.NewFilter()
	if *stockedFlag {
		filter.StockedOnly()
	}
	// if *samplesFlag {
	// filter.SamplesOnly()
	// }
	filter.Types(strings.Split(*teaTypes, ","))

	db, err := hgtealib.New(teas_url, log_url, *proxyStr)
	if err != nil {
		log.Fatal(err)
	}

	command := flag.Arg(0)

	if command == "ls" {
		/*
			// Id            int
			// Name          string
			// Type          string
			// Picked.Year  int
			// Picked.Flush Flush
			// Origin.Country string
			// Origin.Region  string
			Storage.Stocked bool
			Storage.Aging   bool
			Purchased.Location  string
			Purchased.Date      string
			Purchased.Price     float64
			Purchased.Packaging int
			Size          string
			LeafGrade     string
		*/
		headerFmt := "%3s\t%-60s\t%-15s\t%4s\t%9s\t%30s\t%6s\t%6s\t%6s\t%6s\n"
		teaFmt := "%3d\t%-60s\t%-15s\t%d\t%9s\t%30s\t%7d\t%6d\t%6d\t%6d\n"
		if *noPrettyPrintFlag {
			re := regexp.MustCompile("%-?[0-9]+")
			headerFmt = re.ReplaceAllString(headerFmt, "%")
			teaFmt = re.ReplaceAllString(teaFmt, "%")
		}

		fmt.Printf(headerFmt, "Id", "Name", "Type", "Year", "Flush", "Origin", "Entries", "Avg", "Median", "Mode")
		teas, _ := db.Teas(filter)
		for _, tea := range teas {
			fmt.Printf(teaFmt, tea.Id, tea.Name, tea.Type, tea.Picked.Year, tea.Picked.Flush, tea.Origin.String(), tea.LogLen(), tea.Average(), tea.Median(), tea.Mode())
		}
	} else if command == "log" {
		headerFmt := "%-21s\t%-60s\t%10s\t%s\t%s\t%s\t%s\t%s\n"
		entryFmt := "%s\t%-60s\t%10s\t%d\t%v\t%d\t%d\t%s\n"
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
		fmt.Printf(headerFmt, "Time", "Tea", "Steep Time", "Rating", "Fixins", "Vessel", "Temp", "Session")

		// teas, _ := db.Teas(filter)
		log, _ := db.Log(filter)
		for _, v := range log {
			tea, _ := db.Tea(v.Id)
			fmt.Printf(entryFmt, v.DateTime.Format(time.RFC822Z), tea.String(), v.SteepTime, v.Rating, v.Fixins, v.SteepingVessel, v.SteepingTemperature, v.SessionInstance)
		}
	}
}
