package openweatherapi

import "encoding/json"

type weatherUnits string

func (wu weatherUnits) String() string {
	return string(wu)
}

const (
	standardUnits weatherUnits = "standard"
	metricUnits   weatherUnits = "metric"
	imperialUnits weatherUnits = "imperial"
)

type currentWeatherResponse struct {
	Name  string `json:"name"`
	Coord struct {
		Lat float32 `json:"lat"`
		Lon float32 `json:"lon"`
	} `json:"coord"`
	Condition string
	Climate   string
	Temp      float32
}

func (cwr *currentWeatherResponse) GetCondition() (weather string) {
	if cwr != nil {
		weather = cwr.Condition
	}
	return
}

func (cwr *currentWeatherResponse) GetClimate() (climate string) {
	if cwr != nil {
		climate = cwr.Climate
	}
	return
}

func (cwr *currentWeatherResponse) UnmarshalJSON(data []byte) error {
	type alias currentWeatherResponse
	type capture struct {
		*alias
		Weather []struct {
			Main string `json:"main"`
		} `json:"weather"`
		Main struct {
			Temp float32 `json:"temp"`
		} `json:"main"`
	}
	c := &capture{alias: (*alias)(cwr)}
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}
	if len(c.Weather) != 0 {
		cwr.Condition = c.Weather[0].Main
	}
	cwr.Temp = c.Main.Temp
	cwr.Climate = temperatureToClimate(c.Main.Temp)
	return nil
}

func temperatureToClimate(temp float32) string {
	if temp <= climateImperialThresholdCold {
		return climateCold
	} else if temp <= climateImperialThresholdModerate {
		return climateModerate
	}
	return climateHot
}
