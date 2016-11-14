package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/hokiegeek/hgtealib"
	"log"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
	"time"
)

type viewOptions struct {
	delimeter string
	porcelain bool
	fields    []string
}

type options struct {
	Delimeter string              `json:"delimeter"`
	Porcelain bool                `json:"porcelain"`
	Fields    map[string][]string `json:"fields"`
	DbCfg     struct {
		DbType     string `json:"dbType"`
		TeasUrl    string `json:"teasUrl"`
		JournalUrl string `json:"journalUrl"`
	} `json:"dbCfg"`
	Proxy   string           `json:"proxy"`
	filter  *hgtealib.Filter `json:"-"`
	command string           `json:"-"`
}

func newOptions() *options {
	o := new(options)
	o.Delimeter = "\t"

	o.DbCfg.DbType = "tsv"
	o.DbCfg.TeasUrl = "https://docs.google.com/spreadsheets/d/1-U45bMxRE4_n3hKRkTPTWHTkVKC8O3zcSmkjEyYFYOo/pub?output=tsv"
	o.DbCfg.JournalUrl = "https://docs.google.com/spreadsheets/d/1pHXWycR9_luPdHm32Fb2P1Pp7l29Vni3uFH_q3TsdbU/pub?output=tsv"

	o.Fields = make(map[string][]string)
	o.Fields["ls"] = []string{"Id", "Name", "Type", "Year", "Flush", "Origin", "Entries", "Avg", "Median", "Mode"}
	o.Fields["log"] = []string{"Time", "Tea", "Steep Time", "Rating", "Fixins", "Vessel"}
	return o
}

func printHeader(fields map[string]string, opts viewOptions) {
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
			if i == len(opts.fields)-1 {
				fmt.Println()
			}
		}
	}
}

func printTeas(teas map[int]hgtealib.Tea, opts viewOptions) {
	fields := map[string]string{
		"Id":      "%3d",
		"Name":    "%-60s",
		"Type":    "%-15s",
		"Year":    "%d",
		"Flush":   "%9s",
		"Origin":  "%30s",
		"Size":    "%12s",
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
				if tea.Picked.Year == 0 {
					fmt.Printf("%s", "")
				} else {
					fmt.Printf(fields[field], tea.Picked.Year)
				}
			case field == "Flush":
				fmt.Printf(fields[field], tea.Picked.Flush)
			case field == "Origin":
				fmt.Printf(fields[field], tea.Origin.String())
			case field == "Size":
				fmt.Printf(fields[field], tea.Size)
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

func printEntries(db *hgtealib.TeaDb, log []hgtealib.Entry, opts viewOptions) {
	fields := map[string]string{
		"Time":       "%-21s",
		"Tea":        "%-60s",
		"Steep Time": "%10s",
		"Rating":     "%d",
		"Fixins":     "%-25s",
		"Vessel":     "%-15s",
		"Temp":       "%d°",
		"Session":    "%-35s",
		"Comments":   "%s",
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
				var buf bytes.Buffer
				for i, f := range v.Fixins {
					if i != 0 {
						buf.WriteString(", ")
					}
					buf.WriteString(f.String())
				}
				fmt.Printf(fields[field], buf.String())
			case field == "Vessel":
				fmt.Printf(fields[field], v.SteepingVessel)
			case field == "Temp":
				fmt.Printf(fields[field], v.SteepingTemperature)
			case field == "Session":
				fmt.Printf(fields[field], v.SessionInstance)
			case field == "Comments":
				fmt.Printf(fields[field], v.Comments)
			}
		}
		fmt.Println()
	}
}

func parseConfigFile(opts *options, path string) (*options, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Did not find file at: %s", path))
	}

	if err := json.NewDecoder(file).Decode(&opts); err != nil {
		return nil, errors.New(fmt.Sprintf("Encountered error decoding '%s': %s", path, err))
	}

	return opts, nil
}

func parseCommandLineArguments(opts *options) (*options, []string, error) {
	databaseTypeStr := flag.String("dbType", "", "The type of database that the URLs are pointing to")
	proxyStr := flag.String("proxy", "", "Use the given proxy")

	teaTypes := flag.String("types", "", "Comma-delimited list of tea types to select")
	stockedFlag := flag.Bool("stocked", false, "Only display stocked teas")
	// samplesFlag := flag.Bool("samples", false, "Only display tea samples")

	porcelainFlag := flag.Bool("porcelain", false, "Prints out the data in a highly script consumable way")
	fieldsStr := flag.String("fields", "*", "Comma-delimited list of the fields to display")
	// sortStr := flag.String("sort", "*", "Comma-delimited list of fields to sort the display by")

	flag.Parse()

	if *proxyStr != "" {
		opts.Proxy = *proxyStr
	}

	if *databaseTypeStr != "" {
		opts.DbCfg.DbType = *databaseTypeStr
	}

	opts.Porcelain = *porcelainFlag

	opts.filter = hgtealib.NewFilter()
	if *stockedFlag {
		opts.filter.StockedOnly()
	}
	opts.filter.Types(strings.Split(*teaTypes, ","))

	opts.command = flag.Arg(0)

	if *fieldsStr != "*" {
		opts.Fields[opts.command] = strings.Split(*fieldsStr, ",")
	}

	return opts, flag.Args(), nil
}

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	opts := newOptions()
	opts, err = parseConfigFile(opts, path.Join(usr.HomeDir, ".hgteas.json"))
	if err != nil {
		panic(err)
	}
	opts, _, err = parseCommandLineArguments(opts)
	if err != nil {
		panic(err)
	}

	var db *hgtealib.TeaDb
	switch opts.DbCfg.DbType {
	case "tsv":
		db, err = hgtealib.NewFromTsv(opts.DbCfg.TeasUrl, opts.DbCfg.JournalUrl, opts.Proxy)
	}
	if err != nil {
		log.Fatal(err)
	}

	viewOpts := viewOptions{
		delimeter: opts.Delimeter,
		porcelain: opts.Porcelain,
		fields:    opts.Fields[opts.command],
	}

	switch opts.command {
	case "ls":
		teas, _ := db.Teas(opts.filter)
		printTeas(teas, viewOpts)
	case "log":
		log, _ := db.Log(opts.filter)
		printEntries(db, log, viewOpts)
	default:
		log.Fatalf("Unrecognized command: %s\n", opts.command)
	}
}
