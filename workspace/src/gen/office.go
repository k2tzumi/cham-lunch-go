package main

import (
	"fmt"
	"os"

	"github.com/k2tzumi/cham-lunch-go/src/jma"
)

func main() {
	const CITY_NAME = "磐田市"
	var cityName string
	if len(os.Args) < 2 {
		cityName = CITY_NAME
	} else {
		cityName = os.Args[1]
	}

	area, err := jma.GetArea()
	if err != nil {
		panic(err)
	}
	officeCode, err := jma.FindOfficeCode(area, cityName)
	if err != nil {
		panic(err)
	}
	fmt.Println(officeCode)

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
