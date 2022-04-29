package jma

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
