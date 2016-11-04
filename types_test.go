package main

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

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

// Timestamp       Date    Time    Tea     Rating  Comments        Pictures        Steep Time      Steeping Vessel Steep Temperature       Session Instance        Fixins
func TestCreateEntry(t *testing.T) {
	original_entry := []string{time.Now().String(),
		"1/1/2016",
		"1300",
		"42",
		"3",
		"foobar",
		"raboof",
		"5m 20s",
		"0",
		"180",
		"adi9ao3ahda92h",
		"0;7"}

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
	t.Skip("TODO")
}

func TestCreateBadTea(t *testing.T) {
	t.Skip("TODO")
}
