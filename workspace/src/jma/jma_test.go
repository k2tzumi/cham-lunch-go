package jma

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
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

func TestGetArea(t *testing.T) {
	httpmock.Activate()

	f, err := os.Open("../../resources/area.json")
	if err != nil {
		t.Fatalf("resources not found error")
	}
	defer f.Close()

	byteArray, _ := ioutil.ReadAll(f)

	expected := &Area{}
	json.Unmarshal(byteArray, &expected)

	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://www.jma.go.jp/bosai/common/const/area.json",
		httpmock.NewBytesResponder(200, byteArray),
	)

	actual, _ := GetArea()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("want %s, actual %s", expected, actual)
	}
}
