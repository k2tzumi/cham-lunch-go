package jma

type WeatherCode [2]string

const (
	CLEAR WeatherCode = {"100", "晴れ"}
	// PARTLY_CLOUDY WeatherCode = "101"
)

func (w WeatherCode) GetName() string {
	switch w {
	case CLEAR:
		return "晴"
	case PARTLY_CLOUDY:


	}
}
