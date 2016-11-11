package hgtealib

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var testEntries = []Entry{
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

var testTeas = []Tea{
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
		for i := rand.Intn(30); i != 0; i-- {
			e := createRandomEntry()
			t.Add(*e)
		}
	}

	// t.average       int
	// t.median        int
	// t.mode          int

	return t
}

func TestEntryEquality(t *testing.T) {
	if !testEntries[0].Equal(&testEntries[0]) {
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

	err := e.ParseDateTime(fecha, tiempo)
	if err != nil {
		t.Error(err)
	}

	fecha_found := e.DateTime.Format("1/02/2006")
	if fecha_found != fecha {
		t.Fatalf("Expected date to be %s but found %s", fecha, fecha_found)
	}

	tiempo_found := e.DateTime.Format("1504")
	if tiempo_found != tiempo {
		t.Fatalf("Expected time to be %s but found %s", tiempo, tiempo_found)
	}

	// Test for failure
	err = e.ParseDateTime("foo", "bar")
	if err == nil {
		t.Fatal("Incorrectly parsed a string instead of a time value")
	}

	err = e.ParseDateTime("20/50/", tiempo)
	if err == nil {
		t.Fatal("Incorrectly parsed a date time with a badly formatted date")
	}

	err = e.ParseDateTime("10/11/YYYY", tiempo)
	if err == nil {
		t.Fatal("Incorrectly parsed a date time with a text date")
	}

	err = e.ParseDateTime("10/DD/1314", tiempo)
	if err == nil {
		t.Fatal("Incorrectly parsed a date time with a text date")
	}

	err = e.ParseDateTime("MM/11/1314", tiempo)
	if err == nil {
		t.Fatal("Incorrectly parsed a date time with a text date")
	}

	err = e.ParseDateTime(fecha, "12MM")
	if err == nil {
		t.Fatal("Incorrectly parsed a date time with text minutes")
	}

	err = e.ParseDateTime(fecha, "HH34")
	if err == nil {
		t.Fatal("Incorrectly parsed a date time with a text hours")
	}

	// TODO: the ParseDateTime function needs to do some valiation
	// err = e.ParseDateTime("40/50/1", tiempo)
	// if err == nil {
	// 	t.Fatal("Incorrectly parsed a date time with a bogus date")
	// }

	// err = e.ParseDateTime(fecha, "5678")
	// if err == nil {
	// 	t.Fatal("Incorrectly parsed a date time with a bogus time")
	// }

	// err = e.ParseDateTime("40/50/1", "5678")
	// if err == nil {
	// 	t.Fatal("Incorrectly parsed a date time with a bogus date and time")
	// }

	err = e.ParseDateTime(fecha, "13")
	if err == nil {
		t.Fatal("Incorrectly parsed a time without enough digits")
	}

	err = e.ParseDateTime("", tiempo)
	if err == nil {
		t.Fatal("Incorrectly parsed a date time with empty date value")
	}

	err = e.ParseDateTime(fecha, "")
	if err == nil {
		t.Fatal("Incorrectly parsed a date time with empty time value")
	}

	err = e.ParseDateTime("", "")
	if err == nil {
		t.Fatal("Incorrectly parsed a date time with all values blank")
	}
}

func TestEntryParseSteepTime(t *testing.T) {
	e := createRandomEntry()

	// Test good duration
	err := e.ParseSteepTime("4m 20s")
	if err != nil {
		t.Error(err)
	}

	if e.SteepTime != time.Duration(260*1e9) {
		t.Fatal("Steep time was not parsed correctly")
	}

	// Test failure
	err = e.ParseSteepTime("foobar")
	if err == nil {
		t.Fatal("Incorrectly parsed a string instead of a time value")
	}

	err = e.ParseSteepTime("4u 70s")
	if err == nil {
		t.Fatal("Incorrectly parsed a steep time with a bogus unit")
	}

	err = e.ParseSteepTime("")
	if err == nil {
		t.Fatal("Incorrectly parsed an empty value")
	}
}

func TestTeaEquality(t *testing.T) {
	if !testTeas[0].Equal(&testTeas[0]) {
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
