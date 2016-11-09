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

func printTeas(teas map[int]hgtealib.Tea, columns []string) {
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
	delimeter := "\t"
	formats := map[string]string{
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
	}

	// Print the header
	re := regexp.MustCompile("[a-z]+")
	for i, col := range columns {
		if i != 0 {
			fmt.Print(delimeter)
		}
		/*
			if *rawFlag {
				re := regexp.MustCompile("-?[0-9]+")
				formats[col] = re.ReplaceAllString(formats[col], "")
			}
		*/
		fmt.Printf(re.ReplaceAllString(formats[col], "s"), col)
	}
	fmt.Println()

	// Now print the teas
	for _, tea := range teas {
		for i, col := range columns {
			if i != 0 {
				fmt.Print(delimeter)
			}
			switch {
			case col == "Id":
				fmt.Printf(formats[col], tea.Id)
			case col == "Name":
				fmt.Printf(formats[col], tea.Name)
			case col == "Type":
				fmt.Printf(formats[col], tea.Type)
			case col == "Year":
				fmt.Printf(formats[col], tea.Picked.Year)
			case col == "Flush":
				fmt.Printf(formats[col], tea.Picked.Flush)
			case col == "Origin":
				fmt.Printf(formats[col], tea.Origin.String())
			case col == "Entries":
				fmt.Printf(formats[col], tea.LogLen())
			case col == "Avg":
				fmt.Printf(formats[col], tea.Average())
			case col == "Median":
				fmt.Printf(formats[col], tea.Median())
			case col == "Mode":
				fmt.Printf(formats[col], tea.Mode())
			}
		}
		fmt.Println()
	}
}

func printEntries(db *hgtealib.HgTeaDb, log []hgtealib.Entry, columns []string) {
	/*
		// DateTime            time.Time
		// Rating              int
		Comments            string
		// SteepTime           time.Duration
		SteepingVessel      int
		// SteepingTemperature int
		// SessionInstance     string
		Fixins              []string
	*/
	delimeter := "\t"
	formats := map[string]string{
		"Time":       "%-21s",
		"Tea":        "%-60s",
		"Steep Time": "%10s",
		"Rating":     "%d",
		"Fixins":     "%v",
		"Vessel":     "%d",
		"Temp":       "%d°",
		"Session":    "%s",
	}

	// Print the header
	re := regexp.MustCompile("[a-z]+")
	for i, col := range columns {
		if i != 0 {
			fmt.Print(delimeter)
		}
		/*
			if *rawFlag {
				re := regexp.MustCompile("-?[0-9]+")
				formats[col] = re.ReplaceAllString(formats[col], "")
			}
		*/
		fmt.Printf(re.ReplaceAllString(formats[col], "s"), col)
	}
	fmt.Println()

	for _, v := range log {
		tea, _ := db.Tea(v.Tea)
		for i, col := range columns {
			if i != 0 {
				fmt.Print(delimeter)
			}
			switch {
			case col == "Time":
				fmt.Printf(formats[col], v.DateTime.Format(time.RFC822Z))
			case col == "Tea":
				fmt.Printf(formats[col], tea.String())
			case col == "Steep Time":
				fmt.Printf(formats[col], v.SteepTime)
			case col == "Rating":
				fmt.Printf(formats[col], v.Rating)
			case col == "Fixins":
				fmt.Printf(formats[col], v.Fixins)
			case col == "Vessel":
				fmt.Printf(formats[col], v.SteepingVessel)
			case col == "Temp":
				fmt.Printf(formats[col], v.SteepingTemperature)
			case col == "Session":
				fmt.Printf(formats[col], v.SessionInstance)
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

	// rawFlag := flag.Bool("raw", false, "Formats the table prettily")
	// columnsStr := flag.String("columns", "*", "Comma-delimited list of the columns to display")
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

	command := flag.Arg(0)

	if command == "ls" {
		teas, _ := db.Teas(filter)
		printTeas(teas, []string{"Id", "Name", "Type", "Year", "Flush", "Origin", "Entries", "Avg", "Median", "Mode"})
	} else if command == "log" {
		log, _ := db.Log(filter)
		printEntries(db, log, []string{"Time", "Tea", "Steep Time", "Rating", "Fixins", "Vessel", "Temp", "Session"})
	}
}
