package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/k2tzumi/cham-lunch-go/src/jma"
)

//go:embed area.json
var areaJSON []byte

func loadAreaJSON() (*jma.Area, error) {
	area := &jma.Area{}
	if err := json.Unmarshal(areaJSON, &area); err != nil {
		return nil, err
	}
	return area, nil
}

func main() {
	cityName := os.Args[1]
	area, err := loadAreaJSON()
	if err != nil {
		panic("Can not load area.json")
	}
	officeCode, err := area.FindOfficeCode(cityName)
	if err != nil {
		panic("Not found")
	}
	fmt.Println(cityName, officeCode)
	forecasts, err := jma.GetForecast(officeCode)
	if err != nil {
		panic(err)
	}

	fmt.Println(forecasts[1].PublishingOffice)
	fmt.Println(forecasts[1].ReportDatetime.Date())

	for _, forecast := range forecasts {
		if forecast.IsWeekly() {
			fmt.Println(forecast.GetWeeklyWeatherForecast()[0].AreaName)
			fmt.Println(forecast.GetWeeklyWeatherForecast()[0].WeatherCodes)
			fmt.Println(forecast.GetWeeklyTempsForecast()[0].AreaName)
			fmt.Println(forecast.GetWeeklyTempsForecast()[0].TempsMax)
		}
	}
}
