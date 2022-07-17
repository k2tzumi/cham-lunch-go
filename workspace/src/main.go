package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/k2tzumi/cham-lunch-go/src/jma"
	"github.com/slack-go/slack"

	"github.com/aws/aws-lambda-go/lambda"
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
	verificationToken := os.Getenv("VERIFICATION_TOKEN")

	handler := NewHandler(verificationToken)

	http.HandleFunc("/", handler.Handle)

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}

type handler struct {
	verificationToken string
}

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

func NewHandler(verificationToken string) Handler {
	return &handler{verificationToken}
}

func searchWeather(cityName string) (string, error) {
	area, err := loadAreaJSON()
	if err != nil {
		return "", err
	}
	officeCode, err := area.FindOfficeCode(cityName)
	if err != nil {
		return "", err
	}
	forecasts, err := jma.GetForecast(officeCode)
	if err != nil {
		return "", err
	}

	result := []string{}
	for _, forecast := range forecasts {
		if forecast.IsWeekly() {
			result = append(result, fmt.Sprintf(
				"%sの天気。%sが発表、%sの%s時点の予報\n",
				cityName,
				forecasts[1].PublishingOffice,
				forecast.GetWeeklyWeatherForecast()[0].AreaName,
				forecasts[1].ReportDatetime.Format("2006/1/2 15:04"),
			))
			for idx, weatherCode := range forecast.GetWeeklyWeatherForecast()[0].WeatherCodes {
				result = append(result, fmt.Sprintf(
					"%sは%s\n",
					forecast.TimeSeries[0].TimeDefines[idx].Format("1/2"),
					weatherCode.String(),
				))
			}
		}
	}
	return strings.Join(result, "\n"), nil
}

/**
 * Slack用のメッセージを作成する
 */
func buildSlackMessage(text string) (*slack.Msg, error) {
	// TODO: textがhelp or 空のときの処理
	weather, err := searchWeather(text)
	if err != nil {
		// TODO: log.Println で出力するようにする
		return &slack.Msg{
			Text:         fmt.Sprintf("%v", err),
			ResponseType: slack.ResponseTypeInChannel,
		}, nil
	}
	return &slack.Msg{
		Text:         weather,
		ResponseType: slack.ResponseTypeInChannel,
	}, nil
}

func (h *handler) Handle(w http.ResponseWriter, request *http.Request) {
	s, err := slack.SlashCommandParse(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("slash commnd pase error.", err)
		return
	}

	if !s.ValidateToken(h.verificationToken) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("validate token error.")
		return
	}

	params, _ := buildSlackMessage(s.Text)
	b, err := json.Marshal(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("json marshal error.", err)
		return
	}

	// Writerに書き込んでResponseを作る
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		log.Fatal("Response write error.", err)
	}
}
