package hgtealib

import (
	"strings"
	"testing"
)

func TestNewFilter(t *testing.T) {
	if NewFilter() == nil {
		t.Fatal("Cannot create a Filter object")
	}
}

func TestFilterStockedOnly(t *testing.T) {
	if NewFilter().StockedOnly().stockedOnly == false {
		t.Fatal("StockedOnly is not set to true as expected")
	}

	if NewFilter().stockedOnly != false {
		t.Fatal("Default value for StockedOnly is not false")
	}
}

func TestFilterSamplesOnly(t *testing.T) {
	if NewFilter().SamplesOnly().samplesOnly == false {
		t.Fatal("SamplesOnly is not set to true as expected")
	}

	if NewFilter().samplesOnly != false {
		t.Fatal("Default value for SamplesOnly is not false")
	}
}

func TestFilterTypes(t *testing.T) {
	testTypes := []string{"T1", "T2", "T3"}
	f := NewFilter().Types(testTypes)

	for _, v := range testTypes {
		if _, ok := f.types[strings.ToLower(v)]; !ok {
			t.Fatalf("Did not find vilter type '%s' in Filter types map: %+v", v, f.types)
		}
	}

	f = NewFilter().Types([]string{"", "", "", ""})
	if len(f.types) != 0 {
		t.Error("Was able to add empty types to the Filter")
	}
}

func TestFilterType(t *testing.T) {
	testTypes := []string{"T1", "T2", "T3"}
	f := NewFilter()
	for _, v := range testTypes {

		f = f.Type(v)
	}

	for _, v := range testTypes {
		if _, ok := f.types[strings.ToLower(v)]; !ok {
			t.Fatalf("Did not find vilter type '%s' in Filter types map: %+v", v, f.types)
		}
	}

}

func TestNewHgTeaDb(t *testing.T) {
	_, err := newHgTeaDb(testTeas, testEntries)
	if err != nil {
		t.Fatal(err)
	}

	_, err = newHgTeaDb([]*Tea{}, []*Entry{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHgTeaDbLog(t *testing.T) {
	db, err := newHgTeaDb(testTeas, testEntries)
	if err != nil {
		t.Fatal(err)
	}

	log, err := db.Log(NewFilter())
	if err != nil {
		t.Error(err)
	}

	if len(log) != len(testEntries) {
		t.Fatalf("Found %d log entries but expected %d", len(log), len(testEntries))
	}

	for _, l := range log {
		var found bool
		for _, e := range testEntries {
			if l.Equal(e) {
				found = true
			}
		}
		if !found {
			t.Fatal("Did not find log entry in database")
		}
	}
}

func TestHgTeaDbTeas(t *testing.T) {
	db, err := newHgTeaDb(testTeas, testEntries)
	if err != nil {
		t.Fatal(err)
	}

	teas, err := db.Teas(NewFilter())
	if err != nil {
		t.Error(err)
	}

	if len(teas) != len(testTeas) {
		t.Fatalf("Found %d teas but expected %d", len(teas), len(testTeas))
	}

	for _, v := range testTeas {
		if _, ok := teas[v.Id]; !ok {
			t.Fatal("Did not find tea in db")
		}
	}

	for _, tea := range teas {
		for _, l := range tea.Log() {
			var found bool
			for _, e := range testEntries {
				if l.Equal(e) {
					found = true
				}
			}
			if !found {
				t.Fatal("Did not find log entry in tea")
			}
		}
	}
}

func TestHgTeaDbTeasFiltered(t *testing.T) {
	db, err := newHgTeaDb(testTeas, testEntries)
	if err != nil {
		t.Fatal(err)
	}

	// Check that stocked teas are returned
	stockedIds := make([]int, 0)
	for _, v := range testTeas {
		if v.Storage.Stocked {
			stockedIds = append(stockedIds, v.Id)
		}
	}

	filteredTeas, err := db.Teas(NewFilter().StockedOnly())
	if err != nil {
		t.Error(err)
	}

	if len(filteredTeas) != len(stockedIds) {
		t.Fatalf("Expected %d stocked teas but got %d", len(stockedIds), len(filteredTeas))
	}

	for _, id := range stockedIds {
		if _, ok := filteredTeas[id]; !ok {
			t.Fatalf("Expected tea id %d to be in list of stocked teas", id)
		}
	}

	// TODO: do Samples when it is implemented

}

func TestHgTeaDbTea(t *testing.T) {
	db, err := newHgTeaDb(testTeas, testEntries)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range testTeas {
		_, err = db.Tea(v.Id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
