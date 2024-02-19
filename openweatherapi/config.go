package openweatherapi

import "github.com/spf13/viper"

func init() {
	viper.MustBindEnv("open_weather_api_key", "OPEN_WEATHER_API_KEY")

	viper.MustBindEnv("open_weather_units", "OPEN_WEATHER_UNITS")
	viper.SetDefault("open_weather_units", standardUnits)
}

const (
	// open weather api endpoint for current weather, see https://openweathermap.org/current
	openWeatherApiEndpointCurrentWeather string = "https://api.openweathermap.org/data/2.5/weather"

	// open weather api param keys
	openWeatherApiParamAppId string = "appid"
	openWeatherApiParamLat   string = "lat"
	openWeatherApiParamLon   string = "lon"
	openWeatherApiParamUnits string = "units"
)

const (
	// thresholds for climate zones
	// TODO: configurable via environment variable
	climateImperialThresholdCold     = 50
	climateImperialThresholdModerate = 80

	// labels for climate zones
	climateCold     string = "cold"
	climateModerate string = "moderate"
	climateHot      string = "hot"
)
