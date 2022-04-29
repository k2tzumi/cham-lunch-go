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
}
