package hgtealib

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var testTsvTeasHeader = []string{"Timestamp", "Date", "ID", "Name", "Type", "Region", "Year", "Flush", "Purchase Location", "Purchase Date", "Purchase Price", "Ratings", "Comments", "Pictures", "Country", "Leaf Grade", "Blended Teas", "Blend Ratio", "Size", "Stocked", "Aging", "Packaging"}
var testTsvTeas = [][]string{
	[]string{
		time.Now().String(), // 0
		"1/1/2016",
		"42",
		"Name", // 3
		"Type",
		"Region",
		"2016", // 6
		"2",
		"Purchase Location",
		"1/1/2000", // 9
		"99.99",
		"Ratings",
		"Comments", // 12
		"Pictures",
		"Country",
		"Leaf Grade", // 15
		"Blended Teas",
		"Blend Ratio",
		"Size", // 18
		"TRUE",
		"FALSE",
		"0", // 21
	},
}

var testTsvEntriesHeader = []string{"Timestamp", "Date", "Time", "Tea", "Rating", "Comments", "Pictures", "Steep Time", "Steeping Vessel", "Steep Temperature", "Session Instance", "Fixins"}
var testTsvEntries = [][]string{
	[]string{
		time.Now().String(), // 0
		"1/1/2016",
		"1300",
		"42", // 3
		"3",
		"comment",
		"pics", // 6
		"1m 23s",
		"0",
		"180", // 9
		"adi9ao3ahda92h",
		"0;1",
	},
	[]string{
		time.Now().String(), // 0
		"5/12/2010",
		"0300",
		"1", // 3
		"2",
		"comment",
		"pics", // 6
		"4m",
		"0",
		"212", // 9
		"adi9ao3ahda92h",
		"",
	},
	[]string{
		time.Now().String(), // 0
		"1/1/1982",
		"2359",
		"82", // 3
		"1",
		"comment",
		"pics", // 6
		"25s",
		"0",
		"212", // 9
		"74i9ao3ahda92h",
		"7",
	},
	[]string{
		time.Now().String(), // 0
		"1/1/1999",
		"0123",
		"920", // 3
		"1",
		"comment",
		"pics", // 6
		"1m 1s",
		"0",
		"0", // 9
		"01i90o3ahda92h",
		"",
	},
}

func compareTsvArrays(a1, a2 [][]string) error {
	if a1 == nil && a2 == nil {
		return nil
	}

	if a1 == nil || a2 == nil {
		return errors.New("One of the arrays is unexpectedly nil")
	}

	if len(a1) != len(a2) {
		return errors.New(fmt.Sprintf("Found an array of size %d when expected size %d", len(a1), len(a2)))
	}

	for i, row := range a1 {
		if len(row) != len(a2[i]) {
			return errors.New("Incorrect number of fields returned")
		}
		for j, field := range row {
			if field != a2[i][j] {
				return errors.New(fmt.Sprintf("Expected field value of %s at position [%d][%d], but found %s", a2[i][j], i, j, field))
			}
		}
	}

	return nil
}

func CreateTestTea() (*Tea, error) {
	original_tea := testTsvTeas[0]

	tea, err := newTeaFromTsv(original_tea)
	if err != nil {
		return nil, err
	}

	for _, entry := range testTsvEntries {
		e, err := newEntryFromTsv(entry)
		if err != nil {
			return nil, err
		}
		tea.Add(*e)
	}

	return tea, nil
}

func isTsvEqualToEntry(expected []string, received *Entry) (bool, error) {
	// TODO: timestamp
	// if expected[0] != received.Timestamp {

	/*
		// TODO
		if expected[1] != received.Date {
			return false, errors.New(fmt.Sprintf("Date field '%s' did not match expected '%s'", received.Date, expected[1]))
		}

		if expected[2] != received.Time {
			return false, errors.New(fmt.Sprintf("Time field '%s' did not match expected '%s'", received.Time, expected[2]))
		}
	*/

	dummy, _ := strconv.Atoi(expected[3])
	if dummy != received.Tea {
		return false, errors.New(fmt.Sprintf("Tea field %d did not match expected %s", received.Tea, expected[3]))
	}

	dummy, _ = strconv.Atoi(expected[4])
	if dummy != received.Rating {
		return false, errors.New(fmt.Sprintf("Rating field %d did not match expected %s", received.Rating, expected[4]))
	}

	if expected[5] != received.Comments {
		return false, errors.New(fmt.Sprintf("Comments field '%s' did not match expected '%s'", received.Comments, expected[5]))
	}

	// if expected[6] != received.Pictures {
	// return false, errors.New(fmt.Sprintf("SteepTime field '%s' did not match expected '%s'", received.SteepTime, expected[6]))
	// }

	// TODO: compare two durations...
	// if expected[7] != received.SteepTime {
	// return false, errors.New(fmt.Sprintf("SteepTime field '%s' did not match expected '%s'", received.SteepTime, expected[7]))
	// }

	dummy, _ = strconv.Atoi(expected[8])
	if VesselType(dummy) != received.SteepingVessel {
		return false, errors.New(fmt.Sprintf("SteepingVessel field %s did not match expected %s", received.SteepingVessel, expected[8]))
	}

	dummy, _ = strconv.Atoi(expected[9])
	if dummy != received.SteepingTemperature {
		return false, errors.New(fmt.Sprintf("SteepingTemperature field %s did not match expected %s", received.SteepingTemperature, expected[9]))
	}

	if expected[10] != received.SessionInstance {
		return false, errors.New(fmt.Sprintf("SessionInstance field %s did not match expected %s", received.SessionInstance, expected[10]))
	}

	// TODO: if expected[11] != received.Fixins              []string

	return true, nil
}

func isTsvEqualToTea(expected []string, received *Tea) (bool, error) {
	// TODO: timestamp
	// if expected[0] != received.Timestamp {

	// if expected[1] != received.Date {
	// return false, errors.New(fmt.Sprintf("Date field '%s' did not match expected '%s'", received.Date, expected[1]))
	// }

	dummy_int, _ := strconv.Atoi(expected[2])
	if dummy_int != received.Id {
		return false, errors.New(fmt.Sprintf("Id field '%s' did not match expected '%s'", received.Id, expected[2]))
	}

	if expected[3] != received.Name {
		return false, errors.New(fmt.Sprintf("Name field '%s' did not match expected '%s'", received.Name, expected[3]))
	}

	if expected[4] != received.Type {
		return false, errors.New(fmt.Sprintf("Type field '%s' did not match expected '%s'", received.Type, expected[4]))
	}

	if expected[5] != received.Origin.Region {
		return false, errors.New(fmt.Sprintf("Region field '%s' did not match expected '%s'", received.Origin.Region, expected[5]))
	}

	dummy_int, _ = strconv.Atoi(expected[6])
	if dummy_int != received.Picked.Year {
		return false, errors.New(fmt.Sprintf("Year field '%s' did not match expected '%s'", received.Picked.Year, expected[6]))
	}

	dummy_float, _ := strconv.ParseFloat(expected[7], 64)
	if Flush(dummy_float) != received.Picked.Flush {
		return false, errors.New(fmt.Sprintf("Flush field '%s' did not match expected '%s'", received.Picked.Flush, expected[7]))
	}

	if expected[8] != received.Purchased.Location {
		return false, errors.New(fmt.Sprintf("PurchaseLocation field '%s' did not match expected '%s'", received.Purchased.Location, expected[8]))
	}

	if expected[9] != received.Purchased.Date {
		return false, errors.New(fmt.Sprintf("PurchaseDate field '%s' did not match expected '%s'", received.Purchased.Date, expected[9]))
	}

	dummy_float, _ = strconv.ParseFloat(expected[10], 64)
	if dummy_float != received.Purchased.Price {
		return false, errors.New(fmt.Sprintf("PurchasePrice field '%s' did not match expected '%s'", received.Purchased.Price, expected[10]))
	}

	// if expected[11] != received.Ratings {
	// return false, errors.New(fmt.Sprintf("Ratings field '%s' did not match expected '%s'", received.Ratings, expected[11]))
	// }

	// if expected[12] != received.Comments {
	// return false, errors.New(fmt.Sprintf("Comments field '%s' did not match expected '%s'", received.Comments, expected[12]))
	// }

	// if expected[13] != received.Pictures {
	// return false, errors.New(fmt.Sprintf("Pictures field '%s' did not match expected '%s'", received.Pictures, expected[13]))
	// }

	if expected[14] != received.Origin.Country {
		return false, errors.New(fmt.Sprintf("Country field '%s' did not match expected '%s'", received.Origin.Country, expected[14]))
	}

	if expected[15] != received.LeafGrade {
		return false, errors.New(fmt.Sprintf("LeafGrade field '%s' did not match expected '%s'", received.LeafGrade, expected[15]))
	}

	// if expected[16] != received.BlendedTeas {
	// return false, errors.New(fmt.Sprintf("BlendedTeas field '%s' did not match expected '%s'", received.BlendedTeas, expected[16]))
	// }

	// if expected[17] != received.BlendRatio {
	// return false, errors.New(fmt.Sprintf("BlendRatio field '%s' did not match expected '%s'", received.BlendRatio, expected[17]))
	// }

	if expected[18] != received.Size {
		return false, errors.New(fmt.Sprintf("Size field '%s' did not match expected '%s'", received.Size, expected[18]))
	}

	dummy_bool := expected[19] == "TRUE"
	if dummy_bool != received.Storage.Stocked {
		return false, errors.New(fmt.Sprintf("Stocked field '%s' did not match expected '%s'", received.Storage.Stocked, dummy_bool))
	}

	dummy_bool = expected[20] == "TRUE"
	if dummy_bool != received.Storage.Aging {
		return false, errors.New(fmt.Sprintf("Aging field '%s' did not match expected '%s'", received.Storage.Aging, dummy_bool))
	}

	dummy_int, _ = strconv.Atoi(expected[21])
	if dummy_int != received.Purchased.Packaging {
		return false, errors.New(fmt.Sprintf("Packaging field '%s' did not match expected '%s'", received.Purchased.Packaging, expected[21]))
	}

	return true, nil
}

// Timestamp       Date    Time    Tea     Rating  Comments        Pictures        Steep Time      Steeping Vessel Steep Temperature       Session Instance        Fixins
func TestCreateTsvEntry(t *testing.T) {
	original_entry := testTsvEntries[0]

	e, err := newEntryFromTsv(original_entry)
	if err != nil {
		t.Fatalf("Unable to create Entry: %s\n", err)
	}

	if _, err := isTsvEqualToEntry(original_entry, e); err != nil {
		t.Errorf("Entry object does not match expected: %s", err)
	}
}

func TestCreateTsvBadEntry(t *testing.T) {
	if _, err := newEntryFromTsv([]string{time.Now().String(), "TEST"}); err == nil {
		t.Fatal("Successfully created badly formatted entry")
	}
}

func TestCreateTsvTea(t *testing.T) {
	original_tea := testTsvTeas[0]

	tea, err := newTeaFromTsv(original_tea)
	if err != nil {
		t.Fatalf("Unable to create Tea: %s\n", err)
	}

	if _, err := isTsvEqualToTea(original_tea, tea); err != nil {
		t.Errorf("Tea object does not match expected: %s", err)
	}
}

func TestCreateTsvBadTea(t *testing.T) {
	if _, err := newTeaFromTsv([]string{time.Now().String(), "TEST"}); err == nil {
		t.Fatal("Successfully created badly formatted tea")
	}

	bad_tea := make([]string, len(testTsvTeas[0]))
	copy(bad_tea, testTsvTeas[0])
	bad_tea[2] = "one"
	if _, err := newTeaFromTsv(bad_tea); err == nil {
		t.Fatal("Successfully created tea with bad Id")
	}

	copy(bad_tea, testTsvTeas[0])
	bad_tea[6] = "MMXVI"
	if _, err := newTeaFromTsv(bad_tea); err == nil {
		t.Fatal("Successfully created tea with bad Year")
	}

	copy(bad_tea, testTsvTeas[0])
	bad_tea[7] = "FOOBAR"
	if _, err := newTeaFromTsv(bad_tea); err == nil {
		t.Fatal("Successfully created tea with bad Flush")
	}

	copy(bad_tea, testTsvTeas[0])
	bad_tea[10] = "Monsoon"
	if _, err := newTeaFromTsv(bad_tea); err == nil {
		t.Fatal("Successfully created tea with a bad Flush value")
	}

	copy(bad_tea, testTsvTeas[0])
	bad_tea[21] = "Cube"
	if _, err := newTeaFromTsv(bad_tea); err == nil {
		t.Fatal("Successfully created tea with bad Packaging type")
	}
}

func getTsvServer(data [][]string) *httptest.Server {
	// TODO:
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		for i, entry := range data {
			if i != 0 {
				buf.WriteString("\n")
			}
			for j, field := range entry {
				if j != 0 {
					buf.WriteString("\t")
				}
				buf.WriteString(field)
			}
		}

		fmt.Fprint(w, buf.String())
	}))

	return ts
}

func TestGetSheetTsv(t *testing.T) {
	expectedData := [][]string{
		[]string{"T0.1", "T0.2", "T0.3"},
		[]string{"T1.1", "T1.2", ""},
		[]string{"T2.1", "", ""},
	}

	tsvServer := getTsvServer(expectedData)
	defer tsvServer.Close()

	testData, err := getSheetTsv(tsvServer.URL, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := compareTsvArrays(testData, expectedData); err != nil {
		t.Fatal(err)
	}

	// Test for shitty values
	if _, err := getSheetTsv("", ""); err == nil {
		t.Error("Did not receive expected error on blank url")
	}

	if _, err := getSheetTsv("FOOBAR", ""); err == nil {
		t.Error("Did not receive expected error on bad url value")
	}

	tsvBadServer := getTsvServer([][]string{
		[]string{"T0.1", "T0.2"},
		[]string{"T1.1"},
	})
	defer tsvBadServer.Close()

	if _, err := getSheetTsv(tsvBadServer.URL, ""); err == nil {
		t.Error("Did not encounter expected error")
	}

	// _, err = getSheetTsv("http://www.google.com/robots.txt", "")
	// if err != nil {
	// t.Error("Received unexpected error when using random URL that should work")
	// }
}

func TestNewFromTsv(t *testing.T) {
	tsvTeasServer := getTsvServer(append([][]string{testTsvTeasHeader}, testTsvTeas...))
	defer tsvTeasServer.Close()

	tsvEntriesServer := getTsvServer(append([][]string{testTsvEntriesHeader}, testTsvEntries...))
	defer tsvEntriesServer.Close()

	db, err := NewFromTsv(tsvTeasServer.URL, tsvEntriesServer.URL, "")
	if err != nil {
		t.Fatal(err)
	}

	teas, err := db.Teas(NewFilter())
	if err != nil {
		t.Error(err)
	}

	if len(teas) != len(testTsvTeas) {
		t.Fatalf("Expected %d teas but found %d", len(testTsvTeas), len(teas))
	}

	for _, tsv_tea := range testTsvTeas {
		tea_id, _ := strconv.Atoi(tsv_tea[2])
		tea, err := db.Tea(tea_id)
		if err != nil {
			teas, _ := db.Teas(NewFilter())
			t.Fatalf("%s: %v", err, teas)
		}
		if _, err := isTsvEqualToTea(tsv_tea, &tea); err != nil {
			t.Errorf("Tea object does not match expected: %s", err)
		}
	}

	// Test with bad values
	if _, err := NewFromTsv("", "", ""); err == nil {
		t.Error("Did not receive expected error on blank urls")
	}

	if _, err := NewFromTsv("FOO", "BAR", ""); err == nil {
		t.Error("Did not receive expected error on bad url value")
	}
}

func TestNewFromTsvFailure(t *testing.T) {
	noEntriesServer := getTsvServer([][]string{[]string{}})
	defer noEntriesServer.Close()

	if _, err := NewFromTsv(noEntriesServer.URL, noEntriesServer.URL, ""); err == nil {
		t.Error("Did not receive expected error when no data from server")
	}

	tsvBadDataServer := getTsvServer([][]string{
		[]string{"T0.1", "T0.2"},
		[]string{"T1.1"},
	})
	defer tsvBadDataServer.Close()

	tsvIncompleteDataServer := getTsvServer([][]string{
		[]string{"T0.1", "T0.2"},
		[]string{"T1.1", "T1.2"},
	})
	defer tsvIncompleteDataServer.Close()

	tsvEntriesServer := getTsvServer(append([][]string{testTsvEntriesHeader}, testTsvEntries...))
	defer tsvEntriesServer.Close()

	// Should trigger nil when creating the tea struct
	if _, err := NewFromTsv(tsvIncompleteDataServer.URL, tsvEntriesServer.URL, ""); err == nil {
		t.Error("Did not receive expected error when incorrect number of tea fields")
	}

	// Should trigger nil when trying to retrieve data for the journal struct
	if _, err := NewFromTsv(noEntriesServer.URL, "", ""); err == nil {
		t.Error("Did not receive expected error when no valid journal server URL")
	}

	tsvTeasServer := getTsvServer(append([][]string{testTsvTeasHeader}, testTsvTeas...))
	defer tsvTeasServer.Close()

	// Should trigger nil when retrieving journal entries from incomplete data
	if _, err := NewFromTsv(tsvTeasServer.URL, tsvIncompleteDataServer.URL, ""); err == nil {
		t.Error("Did not receive expected error when no data from server")
	}

	// Should trigger nil when retrieving journal entries from an empty server
	if _, err := NewFromTsv(tsvTeasServer.URL, noEntriesServer.URL, ""); err == nil {
		t.Error("Did not receive expected error when no data from server")
	}

	// Should trigger nil when creating the journal struct // FIXME
	if _, err := NewFromTsv(tsvTeasServer.URL, tsvBadDataServer.URL, ""); err == nil {
		t.Error("Did not receive expected error when incorrect number of journal fields")
	}
}
