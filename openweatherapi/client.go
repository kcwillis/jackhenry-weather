package openweatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kcwillis/jackhenry/weather/util"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type client struct {
	logger *zap.Logger
	apiKey string
	units  weatherUnits
}

func newClient() *client {
	return &client{
		apiKey: viper.GetString("open_weather_api_key"),
		units:  weatherUnits(viper.GetString("open_weather_units")),
		logger: util.MustNewLogger(),
	}
}

/*
Write an http server that uses the Open Weather API that exposes an endpoint that takes in lat/long coordinates.
This endpoint should return what the weather condition is outside in that area (snow, rain, etc), whether it’s hot,
cold, or moderate outside (use your own discretion on what temperature equates to each type).
*/
func (c *client) CurrentWeather(ctx context.Context, coord Coord) (CurrentWeatherResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, openWeatherApiEndpointCurrentWeather, nil)
	if err != nil {
		c.logger.Error("failure to create new request", zap.Error(err))
		return nil, err
	}

	q := req.URL.Query()
	q.Set(openWeatherApiParamAppId, c.apiKey)
	q.Set(openWeatherApiParamUnits, c.units.String())
	q.Set(openWeatherApiParamLat, fmt.Sprint(coord.GetLat()))
	q.Set(openWeatherApiParamLon, fmt.Sprint(coord.GetLon()))
	req.URL.RawQuery = q.Encode()

	httpclient := &http.Client{}
	res, err := httpclient.Do(req)
	if err != nil {
		c.logger.Error("failure to send current weather request", zap.String("query", req.URL.RawQuery), zap.Error(err))
		return nil, err
	}

	var response *currentWeatherResponse
	switch res.StatusCode {
	case http.StatusOK:
		response = &currentWeatherResponse{}
		defer res.Body.Close()
		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			c.logger.Error("response body decoder error", zap.Error(err))
			break
		}
	default:
		// TODO: error handling for non-2xx response codes
		resByte, ioerr := io.ReadAll(res.Body)
		if ioerr != nil {
			c.logger.Error("failure to read response body", zap.Error(err))
		}
		err = fmt.Errorf("open weather api call failed with status code: %d, body: %s", res.StatusCode, string(resByte))
	}

	return response, err
}
