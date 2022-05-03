package jma

import (
	_ "embed"
	"encoding/json"
)

type WeatherCode string

//go:embed forecast_const.json
var forecastConstJSON []byte
var forecastConst map[string][]string

func init() {
	if err := json.Unmarshal(forecastConstJSON, &forecastConst); err != nil {
		panic("Can not load forecast_const.json")
	}
}

func (w WeatherCode) GetName() string {
	return forecastConst[string(w)][3]
}
