package jma

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

type WeatherCode string

// Implement Stringer
var _ fmt.Stringer = WeatherCode("")

//go:embed forecast_const.json
var forecastConstJSON []byte
var forecastConst map[string][]string

func init() {
	if err := json.Unmarshal(forecastConstJSON, &forecastConst); err != nil {
		panic("Can not load forecast_const.json")
	}
}

func (w WeatherCode) String() string {
	return forecastConst[string(w)][3]
}
