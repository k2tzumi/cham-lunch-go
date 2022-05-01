package jma

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestFindOfficeCode(t *testing.T) {
	patterns := []struct {
		name          string
		cityName      string
		expected      string
		expectedError bool
	}{
		{"Class20で見つかった場合", "磐田市", "220000", false},
		{"Class20で同一地域で見つかった場合", "浜松市", "220000", false},
		{"Class20で別地域で見つかった場合", "豊田市", "230000", false},
		{"Class15で見つかった場合", "宗谷", "011000", false},
		{"Class10で見つかった場合", "伊豆諸島", "130000", false},
		{"前方一致で見つかること", "奄美", "460040", false},
		{"完全一致で見つかり、結果は同じになること", "奄美地方", "460040", false},
		{"存在しない場合はerrorが返却されること", "unknown", "", true},
	}

	f, err := os.Open("../../resources/area.json")
	if err != nil {
		t.Fatalf("resources not found error")
	}
	defer f.Close()

	byteArray, _ := ioutil.ReadAll(f)

	area := &Area{}
	json.Unmarshal(byteArray, &area)

	for _, pattern := range patterns {
		actual, actualError := FindOfficeCode(area, pattern.cityName)
		if pattern.expected != actual {
			t.Errorf("pattern %s: want %s, actual %s", pattern.name, pattern.expected, actual)
		}
		// FIXME: Check for strict errors
		if (actualError == nil) == pattern.expectedError {
			t.Errorf("pattern %s: want %t, actual %s", pattern.name, pattern.expectedError, actualError)
		}
	}
}
