package server

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/kcwillis/jackhenry/weather/proto/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewWeatherService(t *testing.T) {
	svc := NewWeatherService()
	req := &pb.CurrentWeatherRequest{
		Lat: 51.5073219,
		Lon: -0.1276474,
	}
	res, err := svc.CurrentWeather(context.Background(), req)
	assert.NotNil(t, res)
	assert.NoError(t, err)
	fmt.Printf("res.GetClimate(): %v\n", res.GetClimate())
	fmt.Printf("res.GetCondition(): %v\n", res.GetCondition())
}

func TestNewWeatherServiceValidator(t *testing.T) {
	svc := NewWeatherService()
	{
		req := &pb.CurrentWeatherRequest{
			Lat: 51.5073219,
			Lon: -0.1276474,
		}
		err := svc.validateCurrentWeatherRequest(req)
		assert.NoError(t, err)
	}
	{
		req := &pb.CurrentWeatherRequest{
			Lat: -1111111,
			Lon: -1111111,
		}
		err := svc.validateCurrentWeatherRequest(req)
		assert.Error(t, err)
		fmt.Printf("err: %v\n", err)
		_, err = svc.CurrentWeather(context.Background(), req)
		fmt.Printf("err: %v\n", err)
	}
}
