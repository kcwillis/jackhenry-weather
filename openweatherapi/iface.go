package openweatherapi

import "context"

func NewClient() Client {
	return newClient()
}

/*
Client interface provides interactions with the Open Weather Api
*/
type Client interface {
	CurrentWeather(ctx context.Context, coord Coord) (CurrentWeatherResponse, error)
}

/*
Coord interface is required by the CurrentWeather handler to retrieve lat and lon values
*/
type Coord interface {
	GetLat() float32
	GetLon() float32
}

/*
CurrentWeatherResponse interface provides access to Condition and Climate values
*/
type CurrentWeatherResponse interface {
	GetCondition() string
	GetClimate() string
}
