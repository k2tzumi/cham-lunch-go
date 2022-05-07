package jma

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetForecast(t *testing.T) {
	patterns := []struct {
		publishingOffice string
		officeCode       string
	}{
		{"静岡地方気象台", "220000"},
		{"熊谷地方気象台", "110000"},
	}

	httpmock.Activate()

	for _, pattern := range patterns {
		f, err := os.Open("testdata/forecast/" + pattern.officeCode + ".json")
		if err != nil {
			t.Fatalf("test date not found error")
		}
		defer f.Close()

		byteArray, _ := ioutil.ReadAll(f)

		var forecasts []Forecast
		json.Unmarshal(byteArray, &forecasts)

		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "https://www.jma.go.jp/bosai/forecast/data/forecast/"+pattern.officeCode+".json",
			httpmock.NewBytesResponder(200, byteArray),
		)

		actual, _ := GetForecast(pattern.officeCode)

		assert.Equal(t, pattern.publishingOffice, actual[0].PublishingOffice, pattern.publishingOffice+"の天気予想データが取得できること")
	}
}

func TestGetWeeklyWeatherForecast(t *testing.T) {
	patterns := []struct {
		publishingOffice     string
		officeCode           string
		expectedWeatherCodes []WeatherCode
	}{
		{"静岡地方気象台", "220000", []WeatherCode{"101", "100", "101", "101", "201", "201", "201"}},
		{"熊谷地方気象台", "110000", []WeatherCode{"110", "101", "101", "101", "201", "201", "201"}},
	}

	for _, pattern := range patterns {
		f, err := os.Open("testdata/forecast/" + pattern.officeCode + ".json")
		if err != nil {
			t.Fatalf("test date not found error")
		}
		defer f.Close()

		byteArray, _ := ioutil.ReadAll(f)

		var forecasts []Forecast
		json.Unmarshal(byteArray, &forecasts)

		actual := forecasts[1].GetWeeklyWeatherForecast()

		assert.Equal(t, pattern.expectedWeatherCodes, actual[0].WeatherCodes, pattern.publishingOffice+"の天気コードが取得できること")
	}
}

func TestGetWeeklyTempsForecast(t *testing.T) {
	patterns := []struct {
		publishingOffice string
		officeCode       string
		expectedTempsMin []string
	}{
		{"静岡地方気象台", "220000", []string{"", "11", "10", "12", "14", "16", "16"}},
		{"熊谷地方気象台", "110000", []string{"", "9", "9", "11", "14", "14", "15"}},
	}

	for _, pattern := range patterns {
		f, err := os.Open("testdata/forecast/" + pattern.officeCode + ".json")
		if err != nil {
			t.Fatalf("test date not found error")
		}
		defer f.Close()

		byteArray, _ := ioutil.ReadAll(f)

		var forecasts []Forecast
		json.Unmarshal(byteArray, &forecasts)

		actual := forecasts[1].GetWeeklyTempsForecast()

		assert.Equal(t, pattern.expectedTempsMin, actual[0].TempsMin, pattern.publishingOffice+"の最低気温が取得できること")
	}
}

func TestIsWeekly(t *testing.T) {
	f, err := os.Open("testdata/forecast/220000.json")
	if err != nil {
		t.Fatalf("test date not found error")
	}
	defer f.Close()

	byteArray, _ := ioutil.ReadAll(f)

	var forecasts []Forecast
	json.Unmarshal(byteArray, &forecasts)

	assert.False(t, forecasts[0].IsWeekly(), "天気予想の1件目は明後日までの詳細")
	assert.True(t, forecasts[1].IsWeekly(), "天気予想の2件目は７日先まで")
}
