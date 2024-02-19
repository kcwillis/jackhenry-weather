package openweatherapi

import (
	"context"
	"encoding/json"
	"testing"

	v1 "github.com/kcwillis/jackhenry/weather/proto/v1"
	"github.com/stretchr/testify/assert"
)

func Test_newClient(t *testing.T) {
	c := newClient()
	assert.NotNil(t, c)

	coord := &v1.CurrentWeatherRequest{
		Lat: 51.5073219,
		Lon: -0.1276474,
	}

	response, err := c.CurrentWeather(context.Background(), coord)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func Test_unmarshalResponse(t *testing.T) {
	r := &currentWeatherResponse{}
	err := json.Unmarshal([]byte(respexample), r)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, float32(51.51), r.Coord.Lat)
	assert.Equal(t, float32(-0.13), r.Coord.Lon)
	assert.Equal(t, float32(51.8), r.Temp)
	assert.Equal(t, "Rain", r.Condition)
}

var respexample = `{
	"coord": {
	  "lon": -0.13,
	  "lat": 51.51
	},
	"weather": [
	  {
		"id": 500,
		"main": "Rain",
		"description": "light rain",
		"icon": "10n"
	  }
	],
	"base": "stations",
	"main": {
	  "temp": 51.8,
	  "feels_like": 51.1,
	  "temp_min": 49.78,
	  "temp_max": 53.04,
	  "pressure": 1025,
	  "humidity": 94
	},
	"visibility": 8000,
	"wind": {
	  "speed": 10.36,
	  "deg": 220
	},
	"rain": {
	  "1h": 0.96
	},
	"clouds": {
	  "all": 75
	},
	"dt": 1708223950,
	"sys": {
	  "type": 2,
	  "id": 2075535,
	  "country": "GB",
	  "sunrise": 1708240206,
	  "sunset": 1708276752
	},
	"timezone": 0,
	"id": 2643743,
	"name": "London",
	"cod": 200
  }`
