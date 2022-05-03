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
	forecasts, err := jma.GetForecast(officeCode)
	if err != nil {
		panic(err)
	}

	for _, forecast := range forecasts {
		if forecast.IsWeekly() {
			fmt.Printf(
				"%sの天気。%sが発表、%sの%s時点の予報\n",
				cityName,
				forecasts[1].PublishingOffice,
				forecast.GetWeeklyWeatherForecast()[0].AreaName,
				forecasts[1].ReportDatetime.Format("2006/1/2 15:04"),
			)
			for idx, weatherCode := range forecast.GetWeeklyWeatherForecast()[0].WeatherCodes {
				fmt.Printf(
					"%sは%s ／ 最低気温%s 最高気温%s\n",
					forecast.TimeSeries[0].TimeDefines[idx].Format("1/2"),
					weatherCode.GetName(),
					forecast.GetWeeklyTempsForecast()[0].TempsMin[idx],
					forecast.GetWeeklyTempsForecast()[0].TempsMax[idx],
				)
			}
		}
	}
}
