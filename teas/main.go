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

type formatOpts struct {
	delimeter string
	fields    []string
	porcelain bool
}

func printHeader(fields map[string]string, opts formatOpts) {
	re_lcalpha := regexp.MustCompile("[a-z]+")
	re_dashnums := regexp.MustCompile("-?[0-9]+")
	for i, field := range opts.fields {
		if opts.porcelain {
			fields[field] = re_dashnums.ReplaceAllString(fields[field], "")
		} else {
			if i != 0 {
				fmt.Print(opts.delimeter)
			}
			fmt.Printf(re_lcalpha.ReplaceAllString(fields[field], "s"), field)
			if i == len(opts.fields) {
				fmt.Println()
			}
		}
	}
}

func printTeas(teas map[int]hgtealib.Tea, opts formatOpts) {
	fields := map[string]string{
		"Id":      "%3d",
		"Name":    "%-60s",
		"Type":    "%-15s",
		"Year":    "%d",
		"Flush":   "%9s",
		"Origin":  "%30s",
		"Entries": "%7d",
		"Avg":     "%6d",
		"Median":  "%6d",
		"Mode":    "%6d",
		// Storage.Stocked bool
		// Storage.Aging   bool
		// Purchased.Location  string
		// Purchased.Date      string
		// Purchased.Price     float64
		// Purchased.Packaging int
		// Size          string
		// LeafGrade     string
	}

	printHeader(fields, opts)

	// Now print the teas
	for _, tea := range teas {
		for i, field := range opts.fields {
			if i != 0 {
				fmt.Print(opts.delimeter)
			}
			switch {
			case field == "Id":
				fmt.Printf(fields[field], tea.Id)
			case field == "Name":
				fmt.Printf(fields[field], tea.Name)
			case field == "Type":
				fmt.Printf(fields[field], tea.Type)
			case field == "Year":
				fmt.Printf(fields[field], tea.Picked.Year)
			case field == "Flush":
				fmt.Printf(fields[field], tea.Picked.Flush)
			case field == "Origin":
				fmt.Printf(fields[field], tea.Origin.String())
			case field == "Entries":
				fmt.Printf(fields[field], tea.LogLen())
			case field == "Avg":
				fmt.Printf(fields[field], tea.Average())
			case field == "Median":
				fmt.Printf(fields[field], tea.Median())
			case field == "Mode":
				fmt.Printf(fields[field], tea.Mode())
			}
		}
		fmt.Println()
	}
}

func printEntries(db *hgtealib.HgTeaDb, log []hgtealib.Entry, opts formatOpts) {
	fields := map[string]string{
		"Time":       "%-21s",
		"Tea":        "%-60s",
		"Steep Time": "%10s",
		"Rating":     "%d",
		"Fixins":     "%v",
		"Vessel":     "%d",
		"Temp":       "%d°",
		"Session":    "%s",
		// Comments            string
		// SteepingVessel      int
		// Fixins              []string
	}

	printHeader(fields, opts)

	for _, v := range log {
		tea, _ := db.Tea(v.Tea)
		for i, field := range opts.fields {
			if i != 0 {
				fmt.Print(opts.delimeter)
			}
			switch {
			case field == "Time":
				fmt.Printf(fields[field], v.DateTime.Format(time.RFC822Z))
			case field == "Tea":
				fmt.Printf(fields[field], tea.String())
			case field == "Steep Time":
				fmt.Printf(fields[field], v.SteepTime)
			case field == "Rating":
				fmt.Printf(fields[field], v.Rating)
			case field == "Fixins":
				fmt.Printf(fields[field], v.Fixins)
			case field == "Vessel":
				fmt.Printf(fields[field], v.SteepingVessel)
			case field == "Temp":
				fmt.Printf(fields[field], v.SteepingTemperature)
			case field == "Session":
				fmt.Printf(fields[field], v.SessionInstance)
			}
		}
		fmt.Println()
	}
}

func main() {
	teas_url := "https://docs.google.com/spreadsheets/d/1-U45bMxRE4_n3hKRkTPTWHTkVKC8O3zcSmkjEyYFYOo/pub?output=tsv"
	log_url := "https://docs.google.com/spreadsheets/d/1pHXWycR9_luPdHm32Fb2P1Pp7l29Vni3uFH_q3TsdbU/pub?output=tsv"

	proxyStr := flag.String("proxy", "", "Use the given proxy")

	teaTypes := flag.String("types", "", "Comma-delimited list of tea types to select")
	stockedFlag := flag.Bool("stocked", false, "Only display stocked teas")
	// samplesFlag := flag.Bool("samples", false, "Only display tea samples")

	porcelainFlag := flag.Bool("porcelain", false, "Prints out the data in a highly script consumable way")
	fieldsStr := flag.String("fields", "*", "Comma-delimited list of the fields to display")
	// sortStr := flag.String("sort", "*", "Comma-delimited list of fields to sort the display by")

	flag.Parse()

	filter := hgtealib.NewFilter()
	if *stockedFlag {
		filter.StockedOnly()
	}
	filter.Types(strings.Split(*teaTypes, ","))

	db, err := hgtealib.NewFromTsv(teas_url, log_url, *proxyStr)
	if err != nil {
		log.Fatal(err)
	}

	fields := strings.Split(*fieldsStr, ",")

	command := flag.Arg(0)

	switch {
	case command == "ls":
		if len(fields) == 1 && fields[0] == "*" {
			fields = []string{"Id", "Name", "Type", "Year", "Flush", "Origin", "Entries", "Avg", "Median", "Mode"}
		}
		teas, _ := db.Teas(filter)
		printTeas(teas, formatOpts{delimeter: "\t", fields: fields, porcelain: *porcelainFlag})
	case command == "log":
		if len(fields) == 1 && fields[0] == "*" {
			fields = []string{"Time", "Tea", "Steep Time", "Rating", "Fixins", "Vessel", "Temp", "Session"}
		}
		log, _ := db.Log(filter)
		printEntries(db, log, formatOpts{delimeter: "\t", fields: fields, porcelain: *porcelainFlag})
	}
}
