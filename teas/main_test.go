package main

func ExamplePrintHeader() {
	testFields := map[string]string{
		"T0": "%40s",
		"T1": "%-4s",
		"T2": "%s",
		"T3": "%3d",
	}
	testOpts := formatOpts{
		delimeter: " ",
		fields:    []string{"T1", "T2", "T3"},
		porcelain: false,
	}

	printHeader(testFields, testOpts)

	// Output: T1   T2  T3
}

func ExamplePrintTeas() {
	// fields = []string{"Id", "Name", "Type", "Year", "Flush", "Origin", "Entries", "Avg", "Median", "Mode"}
	// printTeas(TODO)
}

func ExamplePrintEntries() {
	// fields = []string{"Time", "Tea", "Steep Time", "Rating", "Fixins", "Vessel", "Temp", "Session", "Comments"}
	// printEntries(TODO)
}
