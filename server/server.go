package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/kcwillis/jackhenry/weather/proto/v1"
	"github.com/kcwillis/jackhenry/weather/util"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	viper.BindEnv("http_port", "HTTP_PORT")
	viper.SetDefault("http_port", 8090)

	viper.BindEnv("grpc_port", "GRPC_PORT")
	viper.SetDefault("grpc_port", 8080)
}
func Serve(ctx context.Context) error {
	logger := util.GetGlobalLogger()
	grpcPort := viper.GetInt("grpc_port")
	httpPort := viper.GetInt("http_port")
	errChan := make(chan error, 2)

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Error("listener failed", zap.Error(err))
		return err
	}

	// create new weather service
	svc := NewWeatherService()

	// Create a gRPC server object
	s := grpc.NewServer()

	// Attach the Greeter service to the server
	pb.RegisterWeatherServiceServer(s, svc)
	// Serve gRPC server
	go func() {
		logger.Info(fmt.Sprintf("Serving gRPC 0.0.0.0:%d", grpcPort))
		if err := s.Serve(lis); err != nil {
			logger.Error("gRPC server failure", zap.Error(err))
			errChan <- err
		}
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0:%d", grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("Failed to dial server", zap.Error(err))
		return err
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	pb.RegisterWeatherServiceHandler(ctx, gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: gwmux,
	}

	go func() {
		logger.Info(fmt.Sprintf("Serving gRPC-Gateway http://0.0.0.0:%d", httpPort))
		if err := gwServer.ListenAndServe(); err != nil {
			logger.Error("http server failure", zap.Error(err))
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-ctx.Done():
	}

	return nil
}
