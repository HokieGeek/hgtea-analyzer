package hgtealib

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

var testTeas = [][]string{
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

var testEntries = [][]string{
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
		"212", // 9
		"01i90o3ahda92h",
		"",
	},
}

var expectedValues = []int{
	1, // Average
	2, // Median
	1, // Mode
}

func CreateTestTea() (*Tea, error) {
	original_tea := testTeas[0]

	tea, err := newTea(original_tea)
	if err != nil {
		return nil, err
	}

	for _, entry := range testEntries {
		err = tea.Add(entry)
		if err != nil {
			return nil, err
		}
	}

	return tea, nil
}

func AreEntriesEqual(expected []string, received *Entry) (bool, error) {
	// TODO: timestamp
	// if expected[0] != received.Timestamp {

	if expected[1] != received.Date {
		return false, errors.New(fmt.Sprintf("Date field '%s' did not match expected '%s'", received.Date, expected[1]))
	}

	if expected[2] != received.Time {
		return false, errors.New(fmt.Sprintf("Time field '%s' did not match expected '%s'", received.Time, expected[2]))
	}

	// dummy, _ := strconv.Atoi(expected[3])
	// if dummy != received.Tea {
	// 	return false, errors.New(fmt.Sprintf("Tea field %d did not match expected %s", received.Tea, expected[3]))
	// }

	dummy, _ := strconv.Atoi(expected[4])
	if dummy != received.Rating {
		return false, errors.New(fmt.Sprintf("Rating field %d did not match expected %s", received.Rating, expected[4]))
	}

	if expected[5] != received.Comments {
		return false, errors.New(fmt.Sprintf("Comments field '%s' did not match expected '%s'", received.Comments, expected[5]))
	}

	// if expected[6] != received.Pictures {
	// return false, errors.New(fmt.Sprintf("SteepTime field '%s' did not match expected '%s'", received.SteepTime, expected[6]))
	// }

	if expected[7] != received.SteepTime {
		return false, errors.New(fmt.Sprintf("SteepTime field '%s' did not match expected '%s'", received.SteepTime, expected[7]))
	}

	dummy, _ = strconv.Atoi(expected[8])
	if dummy != received.SteepingVessel {
		return false, errors.New(fmt.Sprintf("SteepingVessel field %s did not match expected %s", received.SteepingVessel, expected[8]))
	}

	dummy, _ = strconv.Atoi(expected[9])
	if dummy != received.SteepingTemperature {
		return false, errors.New(fmt.Sprintf("SteepingTemperature field %s did not match expected %s", received.SteepingTemperature, expected[9]))
	}

	if expected[10] != received.SessionInstance {
		return false, errors.New(fmt.Sprintf("SessionInstance field %s did not match expected %s", received.SessionInstance, expected[10]))
	}

	// if expected[11] != received.Fixins              []string

	return true, nil
}

func AreTeasEqual(expected []string, received *Tea) (bool, error) {
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

	// if expected[5] != received.Region {
	// return false, errors.New(fmt.Sprintf("Region field '%s' did not match expected '%s'", received.Region, expected[5]))
	// }

	dummy_int, _ = strconv.Atoi(expected[6])
	if dummy_int != received.Year {
		return false, errors.New(fmt.Sprintf("Year field '%s' did not match expected '%s'", received.Year, expected[6]))
	}

	// dummy_int, _ = strconv.Atoi(expected[7])
	// if dummy_int != received.Flush {
	// return false, errors.New(fmt.Sprintf("Flush field '%s' did not match expected '%s'", received.Flush, expected[7]))
	// }

	// if expected[8] != received.PurchaseLocation {
	// return false, errors.New(fmt.Sprintf("PurchaseLocation field '%s' did not match expected '%s'", received.PurchaseLocation, expected[8]))
	// }

	// if expected[9] != received.PurchaseDate {
	// return false, errors.New(fmt.Sprintf("PurchaseDate field '%s' did not match expected '%s'", received.PurchaseDate, expected[9]))
	// }

	// if expected[10] != received.PurchasePrice {
	// return false, errors.New(fmt.Sprintf("PurchasePrice field '%s' did not match expected '%s'", received.PurchasePrice, expected[10]))
	// }

	// if expected[11] != received.Ratings {
	// return false, errors.New(fmt.Sprintf("Ratings field '%s' did not match expected '%s'", received.Ratings, expected[11]))
	// }

	// if expected[12] != received.Comments {
	// return false, errors.New(fmt.Sprintf("Comments field '%s' did not match expected '%s'", received.Comments, expected[12]))
	// }

	// if expected[13] != received.Pictures {
	// return false, errors.New(fmt.Sprintf("Pictures field '%s' did not match expected '%s'", received.Pictures, expected[13]))
	// }

	// if expected[14] != received.Country {
	// return false, errors.New(fmt.Sprintf("Country field '%s' did not match expected '%s'", received.Country, expected[14]))
	// }

	// if expected[15] != received.LeafGrade {
	// return false, errors.New(fmt.Sprintf("LeafGrade field '%s' did not match expected '%s'", received.LeafGrade, expected[15]))
	// }

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
	if dummy_bool != received.Stocked {
		return false, errors.New(fmt.Sprintf("Stocked field '%s' did not match expected '%s'", received.Stocked, dummy_bool))
	}

	dummy_bool = expected[20] == "TRUE"
	if dummy_bool != received.Aging {
		return false, errors.New(fmt.Sprintf("Aging field '%s' did not match expected '%s'", received.Aging, dummy_bool))
	}

	// dummy_int, _ = strconv.Atoi(expected[21])
	// if dummy_int != received.Packaging {
	// return false, errors.New(fmt.Sprintf("Packaging field '%s' did not match expected '%s'", received.Packaging, expected[21]))
	// }

	return true, nil
}

// Timestamp       Date    Time    Tea     Rating  Comments        Pictures        Steep Time      Steeping Vessel Steep Temperature       Session Instance        Fixins
func TestCreateEntry(t *testing.T) {
	original_entry := testEntries[0]

	e, err := newEntry(original_entry)
	if err != nil {
		t.Fatalf("Unable to create Entry: %s\n", err)
	}

	if _, err := AreEntriesEqual(original_entry, e); err != nil {
		t.Errorf("Entry object does not match expected: %s", err)
	}
}

func TestCreateBadEntry(t *testing.T) {
	incomplete_entry := []string{time.Now().String(), "TEST"}
	_, err := newEntry(incomplete_entry)
	if err == nil {
		t.Fatal("Successfully created badly formatted entry")
	}
}

func TestCreateTea(t *testing.T) {
	original_tea := testTeas[0]

	tea, err := newTea(original_tea)
	if err != nil {
		t.Fatalf("Unable to create Tea: %s\n", err)
	}

	if _, err := AreTeasEqual(original_tea, tea); err != nil {
		t.Errorf("Tea object does not match expected: %s", err)
	}
}

func TestCreateBadTea(t *testing.T) {
	incomplete_tea := []string{time.Now().String(), "TEST"}
	_, err := newTea(incomplete_tea)
	if err == nil {
		t.Fatal("Successfully created badly formatted tea")
	}

	bad_id_tea := make([]string, len(testTeas[0]))
	copy(bad_id_tea, testTeas[0])
	bad_id_tea[2] = "one"
	_, err = newTea(bad_id_tea)
	if err == nil {
		t.Fatal("Successfully created tea with bad Id")
	}

	bad_year_tea := make([]string, len(testTeas[0]))
	copy(bad_year_tea, testTeas[0])
	bad_year_tea[6] = "MMXVI"
	_, err = newTea(bad_year_tea)
	if err == nil {
		t.Fatal("Successfully created tea with bad Year")
	}
}

func TestTeaAdd(t *testing.T) {
	original_tea := testTeas[0]

	tea, err := newTea(original_tea)
	if err != nil {
		t.Fatalf("Unable to create Tea: %s\n", err)
	}

	err = tea.Add(testEntries[0])
	if err != nil {
		t.Fatalf("Error adding entry to tea: %s\n", err)
	}

	// TODO
	// if _, err := AreEntriesEqual(testEntries[0], tea.Log[0]); err != nil {
	// t.Errorf("Entry object does not match expected: %s", err)
	// }

	// This should fail
	err = tea.Add([]string{"TEST"})
	if err == nil {
		t.Fatal("Successfully added bad log entry")
	}
}

func TestTeaAverage(t *testing.T) {
	tea, err := CreateTestTea()
	if err != nil {
		t.Fatalf("Error creating a tea: %s\n", err)
	}

	val := tea.Average()
	if expectedValues[0] != val {
		t.Errorf("Calculated average %d does not match expected: %d\n", val, expectedValues[0])
	}
}

func TestTeaMedian(t *testing.T) {
	tea, err := CreateTestTea()
	if err != nil {
		t.Fatalf("Error creating a tea: %s\n", err)
	}

	val := tea.Median()
	if expectedValues[1] != val {
		t.Errorf("Calculated median %d does not match expected: %d\n", val, expectedValues[1])
	}
}

func TestTeaMode(t *testing.T) {
	tea, err := CreateTestTea()
	if err != nil {
		t.Fatalf("Error creating a tea: %s\n", err)
	}

	val := tea.Mode()
	if expectedValues[2] != val {
		t.Errorf("Calculated mode %d does not match expected: %d\n", val, expectedValues[2])
	}
}

func TestTeaString(t *testing.T) {
	original_tea := testTeas[0]

	tea, err := newTea(original_tea)
	if err != nil {
		t.Fatalf("Unable to create Tea: %s\n", err)
	}

	if len(tea.String()) <= 0 {
		t.Error("Tea String() function returned empty string")
	}
}
