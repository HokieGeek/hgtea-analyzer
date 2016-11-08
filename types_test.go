package hgtealib

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
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
		Fixins:              []string{"Milk", "Sugar"},
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
	e.SteepingVessel = r.Intn(12)
	e.SteepingTemperature = r.Intn(212)
	e.SessionInstance = createRandomString(1)
	e.Fixins = strings.Split(createRandomString(r.Intn(3)), " ")

	return e
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
