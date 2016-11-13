package hgtealib

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

var testEntries = []*Entry{
	{
		Tea:                 42,
		DateTime:            time.Now(),
		Rating:              3,
		Comments:            "These are comments",
		SteepTime:           time.Minute * 3,
		SteepingVessel:      0, // TODO
		SteepingTemperature: 180,
		SessionInstance:     "DEADBEEF",
		Fixins:              []TeaFixin{Milk, Sugar},
	},
}

var testTeas = []*Tea{
	{
		Id:   42,
		Name: "Test Tea #1",
		Type: "Black Flavored",
		Picked: TeaPickPeriod{
			Year:  2009,
			Flush: InBetween,
		},
		Origin: TeaOrigin{
			Country: "India",
			Region:  "Assam",
		},
		Storage: TeaStorageState{
			Stocked: true,
			Aging:   false,
		},
		Purchased: TeaPurchaseInfo{
			Location:  "testing.com",
			Date:      "1/2/2009",
			Price:     1234.56,
			Packaging: 0,
		},
		Size:      "2oz sample",
		LeafGrade: "STFTGFOPOMG!",
		// log           map[time.Time]Entry
		// logSortedKeys TimeSlice
		// average       int
		// median        int
		// mode          int
	},
	{
		Id:   101,
		Name: "Test Tea #2",
		Type: "Black",
		Picked: TeaPickPeriod{
			Year: 2009,
		},
		Origin: TeaOrigin{
			Country: "China",
		},
		Storage: TeaStorageState{
			Stocked: false,
			Aging:   false,
		},
		Purchased: TeaPurchaseInfo{
			Location:  "testing.com",
			Date:      "11/14/2010",
			Price:     19.99,
			Packaging: 0,
		},
		Size:      "2oz",
		LeafGrade: "OP",
		// log           map[time.Time]Entry
		// logSortedKeys TimeSlice
		// average       int
		// median        int
		// mode          int
	},
}

func createRandomString(sentences int) string {
	var buf bytes.Buffer

	var b []byte
	for i := 0; i < sentences; i++ {
		b = make([]byte, 8)
		rand.Read(b)
		buf.WriteString(fmt.Sprintf("%x", b))
	}

	return buf.String()
}

func createRandomEntry() *Entry {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	e := new(Entry)

	e.Tea = r.Int()
	e.DateTime = time.Unix(time.Now().Unix()-r.Int63(), 0)
	e.Rating = r.Intn(4)
	e.Comments = createRandomString(r.Intn(5))
	e.SteepTime = time.Duration(r.Intn(720))
	e.SteepingVessel = VesselType(r.Intn(9))
	e.SteepingTemperature = r.Intn(212)
	e.SessionInstance = createRandomString(1)
	e.Fixins = []TeaFixin{TeaFixin(r.Intn(8)), TeaFixin(r.Intn(8))}

	return e
}

func createRandomTea(withEntries bool) *Tea {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	t := new(Tea)

	t.Id = r.Int()
	t.Name = createRandomString(1)
	t.Type = createRandomString(1)
	t.Picked.Year = r.Int()
	t.Picked.Flush = Flush(r.Intn(5))
	t.Origin.Country = createRandomString(1)
	t.Origin.Region = createRandomString(1)
	t.Storage.Stocked = ((r.Int() % 2) == 0)
	t.Storage.Aging = ((r.Int() % 2) == 0)
	t.Purchased.Location = createRandomString(1)
	t.Purchased.Date = time.Now().Format("1/02/2009")
	t.Purchased.Price = r.Float64()
	t.Purchased.Packaging = r.Intn(10)
	t.Size = createRandomString(1)
	t.LeafGrade = createRandomString(1)

	t.log = make(map[time.Time]Entry)
	t.logSortedKeys = make(TimeSlice, 0)

	if withEntries {
		var numEntries int
		for {
			numEntries = rand.Intn(30)
			if numEntries > 0 {
				break
			}
		}
		for i := numEntries; i != 0; i-- {
			e := createRandomEntry()
			t.Add(*e)
		}
	}

	// t.average       int
	// t.median        int
	// t.mode          int

	return t
}

func TestFlushString(t *testing.T) {
	for _, v := range []Flush{First, InBetween, Second, Monsoon, Autumn} {
		if v.String() == "" {
			t.Error("Flush type did not return a useful string")
		}
	}

	if v := Flush(1.1).String(); v != "" {
		t.Errorf("Expected empty string but found '%s' instead", v)
	}
}

func TestVesselTypeString(t *testing.T) {
	for _, v := range []VesselType{FrenchPress, ShipiaoYixing, TeazerTumbler, TeaStick, MeshSpoon, SaucePan, Cup, Bowl, Gaiwan, Other} {
		if v.String() == "" {
			t.Error("VesselType type did not return a useful string")
		}
	}

	if v := VesselType(-1).String(); v != "" {
		t.Errorf("Expected empty string but found '%s' instead", v)
	}
}

func TestTeaFixinString(t *testing.T) {
	for _, v := range []TeaFixin{Milk, Cream, HalfAndHalf, Sugar, BrownSugar, RawSugar, Honey, VanillaExtract, VanillaBean} {
		if v.String() == "" {
			t.Error("TeaFixin type did not return a useful string")
		}
	}

	if v := TeaFixin(-1).String(); v != "" {
		t.Errorf("Expected empty string but found '%s' instead", v)
	}
}

func TestTeaOriginString(t *testing.T) {
	expected := fmt.Sprintf("%s, %s", testTeas[0].Origin.Region, testTeas[0].Origin.Country)
	if testTeas[0].Origin.String() != expected {
		t.Errorf("Expected origin string '%s' but found '%s'", expected, testTeas[0].Origin.String())
	}

	if testTeas[1].Origin.String() != testTeas[1].Origin.Country {
		t.Errorf("Expected origin string '%s' but found '%s'", testTeas[1].Origin.Country, testTeas[0].Origin.String())
	}
}

func TestEntryEquality(t *testing.T) {
	if !testEntries[0].Equal(testEntries[0]) {
		t.Error("Entry equality identity test failed")
	}

	if createRandomEntry().Equal(createRandomEntry()) {
		t.Error("Entry equality test with random data failed")
	}
}

func TestEntryParseDateTime(t *testing.T) {
	e := createRandomEntry()

	// Test with a good time (rawr)
	fecha := "10/11/1314"
	tiempo := "1234"

	if err := e.ParseDateTime(fecha, tiempo); err != nil {
		t.Error(err)
	}

	if fecha_found := e.DateTime.Format("1/02/2006"); fecha_found != fecha {
		t.Fatalf("Expected date to be %s but found %s", fecha, fecha_found)
	}

	if tiempo_found := e.DateTime.Format("1504"); tiempo_found != tiempo {
		t.Fatalf("Expected time to be %s but found %s", tiempo, tiempo_found)
	}

	// Test for failure
	if e.ParseDateTime("foo", "bar") == nil {
		t.Fatal("Incorrectly parsed a string instead of a time value")
	}

	if e.ParseDateTime("20/50/", tiempo) == nil {
		t.Fatal("Incorrectly parsed a date time with a badly formatted date")
	}

	if e.ParseDateTime("10/11/YYYY", tiempo) == nil {
		t.Fatal("Incorrectly parsed a date time with a text date")
	}

	if e.ParseDateTime("10/DD/1314", tiempo) == nil {
		t.Fatal("Incorrectly parsed a date time with a text date")
	}

	if e.ParseDateTime("MM/11/1314", tiempo) == nil {
		t.Fatal("Incorrectly parsed a date time with a text date")
	}

	if e.ParseDateTime(fecha, "12MM") == nil {
		t.Fatal("Incorrectly parsed a date time with text minutes")
	}

	if e.ParseDateTime(fecha, "HH34") == nil {
		t.Fatal("Incorrectly parsed a date time with a text hours")
	}

	// TODO: the ParseDateTime function needs to do some valiation
	// if e.ParseDateTime("40/50/1", tiempo) == nil {
	// 	t.Fatal("Incorrectly parsed a date time with a bogus date")
	// }

	// if e.ParseDateTime(fecha, "5678") == nil {
	// 	t.Fatal("Incorrectly parsed a date time with a bogus time")
	// }

	// if e.ParseDateTime("40/50/1", "5678") == nil {
	// 	t.Fatal("Incorrectly parsed a date time with a bogus date and time")
	// }

	if e.ParseDateTime(fecha, "13") == nil {
		t.Fatal("Incorrectly parsed a time without enough digits")
	}

	if e.ParseDateTime("", tiempo) == nil {
		t.Fatal("Incorrectly parsed a date time with empty date value")
	}

	if e.ParseDateTime(fecha, "") == nil {
		t.Fatal("Incorrectly parsed a date time with empty time value")
	}

	if e.ParseDateTime("", "") == nil {
		t.Fatal("Incorrectly parsed a date time with all values blank")
	}
}

func TestEntryParseSteepTime(t *testing.T) {
	e := createRandomEntry()

	// Test good duration
	if err := e.ParseSteepTime("4m 20s"); err != nil {
		t.Error(err)
	}

	if e.SteepTime != time.Duration(260*1e9) {
		t.Fatal("Steep time was not parsed correctly")
	}

	// Test failure
	if e.ParseSteepTime("foobar") == nil {
		t.Fatal("Incorrectly parsed a string instead of a time value")
	}

	if e.ParseSteepTime("4u 70s") == nil {
		t.Fatal("Incorrectly parsed a steep time with a bogus unit")
	}

	if e.ParseSteepTime("") == nil {
		t.Fatal("Incorrectly parsed an empty value")
	}
}

func TestTeaEquality(t *testing.T) {
	if !testTeas[0].Equal(testTeas[0]) {
		t.Error("Tea equality identity test failed")
	}

	if createRandomTea(false).Equal(createRandomTea(false)) {
		t.Error("Tea equality test with random data failed")
	}
}

func TestLog(t *testing.T) {
	tea := createRandomTea(false)

	entries := make([]*Entry, rand.Intn(30))
	for i, _ := range entries {
		entries[i] = createRandomEntry()
		tea.Add(*entries[i])
	}

	if tea.LogLen() != len(entries) {
		t.Fatal("LogLen() did not report expected number of entries")
	}

	log := tea.Log()

	if len(log) != tea.LogLen() {
		t.Fatalf("LogLen() and the log do not match ins size")
	}

	// TODO
	// for i,e := range log {
	// }
}

func TestTeaAdd(t *testing.T) {
	tea := createRandomTea(false)
	entry := createRandomEntry()
	tea.Add(*entry)

	log := tea.Log()
	if len(log) != 1 {
		t.Fatalf("Found %d entries when expected 1: %v", len(log), log)
	}

	if !entry.Equal(&log[0]) {
		t.Fatal("Added entry did not match expected")
	}
}

func TestTeaAverage(t *testing.T) {
	for i := rand.Intn(10); i >= 0; i-- {
		var total int

		tea := createRandomTea(true)

		for _, entry := range tea.Log() {
			total += entry.Rating
		}

		avg := total / tea.LogLen()

		if avg != tea.Average() {
			t.Fatalf("Expected average of %d and found %d", avg, tea.Average())
		}
	}
}

func TestTeaMedian(t *testing.T) {
	for i := rand.Intn(10); i >= 0; i-- {
		tea := createRandomTea(true)

		ratings := make([]int, tea.LogLen())
		for i, entry := range tea.Log() {
			ratings[i] = entry.Rating
		}
		sort.Ints(ratings)

		var median int
		if (len(ratings) % 2) == 0 {
			median = (ratings[len(ratings)/2] + ratings[(len(ratings)/2)-1]) / 2
		} else {
			median = ratings[len(ratings)/2]
		}

		if median != tea.Median() {
			t.Logf("Len: %d, ratings: %v", len(ratings), ratings)

			t.Fatalf("Expected median of %d and found %d", median, tea.Median())
		}
	}
}

func TestTeaMode(t *testing.T) {
	for i := rand.Intn(10); i >= 0; i-- {
		tea := createRandomTea(true)

		ratings := make([]int, 5)
		for _, entry := range tea.log {
			ratings[entry.Rating]++
		}

		var max int
		for rating, count := range ratings {
			if count > ratings[max] {
				max = rating
			}
		}

		if max != tea.Mode() {
			t.Errorf("Calculated mode %d does not match expected: %d\n", max, tea.Mode())
		}
	}
}

func TestTeaString(t *testing.T) {
	if createRandomTea(false).String() == "" {
		t.Error("Tea String() function returned empty string")
	}
}
