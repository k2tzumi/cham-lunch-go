package jma

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// https://www.jma.go.jp/bosai/common/const/area.json

type Area struct {
	Offices  map[string]Office  `json:"offices"`
	Class10s map[string]Class10 `json:"class10s"`
	Class15s map[string]Class15 `json:"class15s"`
	Class20s map[string]Class20 `json:"class20s"`
}

type Office struct {
	Name       string `json:"name"`
	EnName     string `json:"enName"`
	Parent     string `json:"parent"`
	OfficeName string `json:"officeName"`
}

type Class10 struct {
	Name   string `json:"name"`
	EnName string `json:"enName"`
	Parent string `json:"parent"`
}

type Class15 struct {
	Name   string `json:"name"`
	EnName string `json:"enName"`
	Parent string `json:"parent"`
}

type Class20 struct {
	Name   string `json:"name"`
	EnName string `json:"enName"`
	Kana   string `json:"kana"`
	Parent string `json:"parent"`
}

func FindOfficeCode(area *Area, cityName string) (string, error) {
	// 指定された市町村名を探索し、areaCodeに突っ込む
	for code := range area.Class20s {
		if strings.HasPrefix(area.Class20s[code].Name, cityName) {
			area20sCode := area.Class20s[code].Parent
			class10sCode := area.Class15s[area20sCode].Parent
			officeCode := area.Class10s[class10sCode].Parent
			return officeCode, nil
		}
	}

	// 見つからなかった場合はエリアを広げてareaCodeを探しに行く
	for code := range area.Class15s {
		if strings.HasPrefix(area.Class15s[code].Name, cityName) {
			class10sCode := area.Class15s[code].Parent
			officeCode := area.Class10s[class10sCode].Parent
			return officeCode, nil
		}
	}

	// 更に見つからなかった場合はエリアを広げてareaCodeを探しに行く
	for code := range area.Class10s {
		if strings.HasPrefix(area.Class10s[code].Name, cityName) {
			officeCode := area.Class10s[code].Parent
			return officeCode, nil
		}
	}

	return "", fmt.Errorf("not found")
}

func GetArea() (*Area, error) {
	url := "https://www.jma.go.jp/bosai/common/const/area.json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	area := &Area{}
	json.Unmarshal(byteArray, &area)

	return area, nil
}
