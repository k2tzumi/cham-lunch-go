package jma

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExtendFile struct {
	*excelize.File
}

func CreateExcel(forecast *Forecast, templatePath string) error {
	f, err := OpenExtendFile(templatePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// areaCode をシート名にする
	areaCode := forecast.TimeSeries[0].WeatherAreas()[0].AreaCode

	// templateシートをバックアップ
	backupName := "origin_backup"
	bakupIndex := f.NewSheet(backupName)
	sheet := "template"
	if err := f.CopySheet(f.GetSheetIndex(sheet), bakupIndex); err != nil {
		return err
	}

	// DefinedNameにセットするvalues
	weeklyForecast := forecast.TimeSeries[0]
	areaName := weeklyForecast.WeatherAreas()[0].AreaName
	reportDatetime := forecast.ReportDatetime.Format("2006年1月2日15時")
	publishingOffice := forecast.PublishingOffice
	// 一つにセルに書き込む値
	cellValues := map[string]string{
		"title":            areaName + "の天気予報（６日先まで）",
		"publishing":       reportDatetime + " " + publishingOffice + " 発表",
		"areaName":         weeklyForecast.WeatherAreas()[0].AreaName,
		"reportDatetime":   reportDatetime,
		"publishingOffice": forecast.PublishingOffice,
		// "averageArea":      forecast.TempAverage.AreaAverages.AreaName,
		// "tempAverageMax":   forecast.TempAverage.AreaAverages.Max,
		// "tempAverageMin":   forecast.TempAverage.AreaAverages.Min,
	}

	// 複数列のセルに書き込む値
	columnsStringValues := map[string]interface{}{
		"days":          weeklyForecast.TimeDefines,
		"weathers":      forecast.GetWeeklyWeatherForecast()[0].WeatherCodes,
		"pops":          forecast.GetWeeklyWeatherForecast()[0].Pops,
		"reliabilities": forecast.GetWeeklyWeatherForecast()[0].Reliabilities,
		"tempsMaxs":     forecast.GetWeeklyTempsForecast()[0].TempsMax,
		"tempsMins":     forecast.GetWeeklyTempsForecast()[0].TempsMin,
	}

	for _, definedName := range f.GetDefinedName() {
		cellRangeAddress := NewCellRangeAddress(definedName.RefersTo)
		// 一つにセルに書き込む
		if cellValue, ok := cellValues[definedName.Name]; ok {
			if err = f.SetCellRefValue(cellRangeAddress.FirstCell, cellValue); err != nil {
				fmt.Println(err)
				return err
			}
			continue
		}
		// 複数列のセルに書き込む
		if columnsValue, ok := columnsStringValues[definedName.Name]; ok {
			if err := f.setSheetColoums(cellRangeAddress.FirstCell, convertAbstractSlice(columnsValue)); err != nil {
				return err
			}
		}
	}

	// エリアコードでシート名を新規作成
	newIndex := f.NewSheet(areaCode)
	// シートをコピーする
	if err := f.CopySheet(f.GetSheetIndex(sheet), newIndex); err != nil {
		return err
	}
	// バックアップから戻す
	if err := f.CopySheet(bakupIndex, f.GetSheetIndex(sheet)); err != nil {
		return err
	}
	// バックアップを削除する
	f.DeleteSheet(backupName)

	// 新規シートをactive化
	f.SetActiveSheet(newIndex)

	// ワークシートを別名保存
	if err = f.SaveAs(areaCode + ".xlsx"); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

type CellRangeAddress struct {
	FirstCell *CellReference
	LastCell  *CellReference
}

type CellReference struct {
	Sheet string
	Row   int
	Col   int
}

func NewCellRangeAddress(refersTo string) *CellRangeAddress {
	r := new(CellRangeAddress)
	splitRefers := strings.Split(refersTo, "!")
	sheet := splitRefers[0]
	axiss := strings.Split(splitRefers[1], ":")
	r.FirstCell = NewCellReference(sheet, axiss[0])
	if len(axiss) == 2 {
		r.LastCell = NewCellReference(sheet, axiss[1])
	}

	return r
}

func NewCellReference(sheet string, axis string) *CellReference {
	c := new(CellReference)
	var err error
	var columnName string
	if columnName, c.Row, err = excelize.SplitCellName(axis); err != nil {
		return nil
	}
	if c.Col, err = excelize.ColumnNameToNumber(columnName); err != nil {
		return nil
	}
	c.Sheet = sheet
	return c
}

func (c *CellReference) GetCellName() (string, error) {
	return excelize.CoordinatesToCellName(c.Col, c.Row)
}

func (c *CellReference) Shift(col int, row int) *CellReference {
	return &CellReference{
		Sheet: c.Sheet,
		Col:   c.Col + col,
		Row:   c.Row + row,
	}
}

func OpenExtendFile(filename string, opt ...excelize.Options) (*ExtendFile, error) {
	f, err := excelize.OpenFile(filename, opt...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &ExtendFile{f}, nil
}

func (f ExtendFile) SetCellRefValue(cellReference *CellReference, value interface{}) error {
	cellName, err := cellReference.GetCellName()
	if err != nil {
		return err
	}

	return f.SetCellValue(cellReference.Sheet, cellName, value)
}

func (f ExtendFile) setSheetColoums(cellReference *CellReference, cellValues []interface{}) error {
	for idx, cellValue := range cellValues {
		// 右に展開
		shiftCell := cellReference.Shift(idx, 0)
		if err := f.SetCellRefValue(shiftCell, cellValue); err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func convertAbstractSlice(array interface{}) []interface{} {
	abstractSlice := []interface{}{}
	switch slice := array.(type) {
	case []string:
		for _, str := range slice {
			abstractSlice = append(abstractSlice, str)
		}
	case []int:
		for _, i := range slice {
			abstractSlice = append(abstractSlice, i)
		}
	case []WeatherCode:
		for _, w := range slice {
			abstractSlice = append(abstractSlice, w)
		}
	case []*time.Time:
		for _, t := range slice {
			abstractSlice = append(abstractSlice, *t)
		}
	}
	return abstractSlice
}
