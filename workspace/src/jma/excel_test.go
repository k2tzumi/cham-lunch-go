package jma

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateExcel(t *testing.T) {
	f, err := os.Open("testdata/forecast/220000.json")
	if err != nil {
		t.Fatalf("test date not found error")
	}
	defer f.Close()

	byteArray, _ := ioutil.ReadAll(f)

	var forecasts []Forecast
	json.Unmarshal(byteArray, &forecasts)

	if err = CreateExcel(&forecasts[1], "../../template.xlsx"); err != nil {
		t.Fatalf("create faild")
	}
}
