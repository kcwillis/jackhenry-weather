package server

import (
	"context"

	validator "github.com/go-playground/validator/v10"
	"github.com/kcwillis/jackhenry/weather/openweatherapi"
	pb "github.com/kcwillis/jackhenry/weather/proto/v1"
	"github.com/kcwillis/jackhenry/weather/util"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewWeatherService() *weatherService {
	return &weatherService{
		openWeatherClient: openweatherapi.NewClient(),
		validate:          validator.New(),
		logger:            util.MustNewLogger(),
	}
}

/*
weatherService is the underlying type which implements the grpc server interface for pb.WeatherServiceServer
*/
type weatherService struct {
	// ensure forward compatibility
	pb.UnimplementedWeatherServiceServer

	logger *zap.Logger

	// a validator to sanitize incoming requests
	validate *validator.Validate

	// a client for interfacing with the open weather api
	openWeatherClient openweatherapi.Client
}

/*
CurrentWeather handler for current weather requests.
Performs validation, and relays request parameters to open weather client which performs api call.
*/
func (s *weatherService) CurrentWeather(ctx context.Context, req *pb.CurrentWeatherRequest) (*pb.CurrentWeatherResponse, error) {
	s.logger.Info("received request", zap.String("endpoint", "CurrentWeather"), zap.Float32("lat", req.GetLat()), zap.Float32("lon", req.GetLon()))

	// validate incoming requests
	if err := s.validateCurrentWeatherRequest(req); err != nil {
		s.logger.Info("invalid request", zap.NamedError("validation_error", err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// relay coordinates to open weather api client
	response, err := s.openWeatherClient.CurrentWeather(ctx, req)
	if err != nil {
		s.logger.Error("open weather api client error", zap.Any("req", req), zap.Error(err))
		return nil, err
	}

	// convert response
	res := &pb.CurrentWeatherResponse{
		Condition: response.GetCondition(),
		Climate:   response.GetClimate(),
	}
	return res, nil
}

/*
validateCurrentWeatherRequest is used to validate the input parameters captured by server
*/
func (s *weatherService) validateCurrentWeatherRequest(req *pb.CurrentWeatherRequest) error {
	/*
		convert request to anonymous struct type with field tags to be used by the validator

		TODO: refactor to use proto plugin which adds custom tags to golang struct fields
		and add 'validate' tags directly to fields in *pb.CurrentWeatherRequest type
	*/
	return s.validate.Struct(struct {
		Lat float32 `validate:"latitude"`
		Lon float32 `validate:"longitude"`
	}{
		Lat: req.GetLat(),
		Lon: req.GetLon(),
	})
}
