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
	t.Skip("TODO")

	// typesFilter := make(map[string]struct{})
	// if len(*teaTypes) > 0 {
	// for _, typeFilter := range strings.Split(*teaTypes, ",") {
	// typesFilter[strings.ToLower(typeFilter)] = struct{}{}
	// }
	// }

	// db, err := buildDatabase(*stockedFlag, *samplesFlag, typesFilter)
	// if err != nil {
	// log.Fatal(err)
	// }

	// fmt.Printf("%-60s %6s %6s %6s %6s\n", "Name", "Num", "Avg", "Median", "Mode")
	// for _, tea := range db {
	// fmt.Printf("%-60s %6d %6d %6d %6d\n", tea.String(), len(tea.Log), tea.Average(), tea.Median(), tea.Mode())
	// }
}

func TestHgTeaDbTeas(t *testing.T) {
	t.Skip("TODO")
}

func TestHgTeaDbTea(t *testing.T) {
	t.Skip("TODO")
}

func TestHgTeaDbLog(t *testing.T) {
	t.Skip("TODO")
}
