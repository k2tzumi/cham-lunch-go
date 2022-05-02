package jma

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// https://www.jma.go.jp/bosai/forecast/data/forecast/110000.json

type Forecast struct {
	PublishingOffice string           `json:"publishingOffice"`
	ReportDatetime   *time.Time       `json:"reportDatetime"`
	TimeSeries       []*TimeSeries    `json:"timeSeries"`
	TempAverage      *json.RawMessage `json:"tempAverage,omitempty"`
	PrecipAverage    *json.RawMessage `json:"precipAverage,omitempty"`
}

type TimeSeries struct {
	TimeDefines []*time.Time       `json:"timeDefines"`
	RawAreas    []*json.RawMessage `json:"areas"`
}

type AreaWeatherForecast struct {
	AreaInfo      `json:"area"`
	WeatherCodes  []WeatherCode `json:"weatherCodes,omitempty"`
	Pops          []string      `json:"pops,omitempty"`
	Reliabilities []string      `json:"reliabilities,omitempty"`
}

type AreaTempsForecast struct {
	AreaInfo      `json:"area"`
	TempsMin      []string `json:"tempsMin,omitempty"`
	TempsMinUpper []string `json:"tempsMinUpper,omitempty"`
	TempsMinLower []string `json:"tempsMinLower,omitempty"`
	TempsMax      []string `json:"tempsMax,omitempty"`
	TempsMaxUpper []string `json:"tempsMaxUpper,omitempty"`
	TempsMaxLower []string `json:"tempsMaxLower,omitempty"`
}

type AreaInfo struct {
	AreaName string `json:"name"`
	AreaCode string `json:"code"`
}

type Average struct {
	AreaAverages *AreaAverage `json:"areas,omitempty"`
}

type AreaAverage struct {
	AreaInfo `json:"area,omitempty"`
	Min      string `json:"min,omitempty"`
	Max      string `json:"max,omitempty"`
}

func (f *Forecast) IsWeekly() bool {
	return f.TempAverage != nil
}

func (t *TimeSeries) IsWeekly() bool {
	return len(t.TimeDefines) == 7
}

func (t *TimeSeries) WeatherAreas() []AreaWeatherForecast {
	var areas []AreaWeatherForecast

	for _, area := range t.RawAreas {
		weather := AreaWeatherForecast{}

		if err := json.Unmarshal(*area, &weather); err != nil {
			return nil
		}
		if len(weather.WeatherCodes) == 0 {
			continue
		}

		areas = append(areas, weather)
	}

	return areas
}

func (t *TimeSeries) TempsAreas() []AreaTempsForecast {
	var areas []AreaTempsForecast

	for _, area := range t.RawAreas {
		temps := AreaTempsForecast{}

		if err := json.Unmarshal(*area, &temps); err != nil {
			return nil
		}
		if len(temps.TempsMin) == 0 {
			continue
		}

		areas = append(areas, temps)
	}

	return areas
}

func (f *Forecast) GetWeeklyWeatherForecast() []AreaWeatherForecast {
	for _, timeSeries := range f.TimeSeries {
		if timeSeries.IsWeekly() {
			weatherAreas := timeSeries.WeatherAreas()

			if weatherAreas != nil {
				return weatherAreas
			}
		}
	}

	return nil
}

func (f *Forecast) GetWeeklyTempsForecast() []AreaTempsForecast {
	for _, timeSeries := range f.TimeSeries {
		if timeSeries.IsWeekly() {
			tempsAreas := timeSeries.TempsAreas()
			if tempsAreas != nil {
				return tempsAreas
			}
		}
	}

	return nil
}

// GetForecast は指定された officeCode の天気予想を取得する
func GetForecast(officeCode string) ([]Forecast, error) {
	url := "https://www.jma.go.jp/bosai/forecast/data/forecast/" + officeCode + ".json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	var forecasts []Forecast
	if err := json.Unmarshal(byteArray, &forecasts); err != nil {
		return forecasts, err
	}

	return forecasts, nil
}
