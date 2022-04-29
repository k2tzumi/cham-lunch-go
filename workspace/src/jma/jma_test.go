package jma

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestFindOfficeCode(t *testing.T) {
	patterns := []struct {
		cityName      string
		expected      string
		expectedError bool
	}{
		{"磐田市", "220000", false},
		{"豊田市", "230000", false},
		{"最上", "060000", false},
		{"東葛飾", "120000", false},
		{"奄美", "460040", false},
		{"unknown", "", true},
	}

	f, err := os.Open("../../resources/area.json")
	if err != nil {
		t.Fatalf("resources not found error")
	}
	defer f.Close()

	byteArray, _ := ioutil.ReadAll(f)

	area := &Area{}
	json.Unmarshal(byteArray, &area)

	for idx, pattern := range patterns {
		actual, actualError := FindOfficeCode(area, pattern.cityName)
		if pattern.expected != actual {
			t.Errorf("pattern %d: want %s, actual %s", idx, pattern.expected, actual)
		}
		// FIXME: Check for strict errors
		if (actualError == nil) == pattern.expectedError {
			t.Errorf("pattern %d: want %t, actual %s", idx, pattern.expectedError, actualError)
		}
	}
}
