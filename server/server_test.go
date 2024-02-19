package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
API server listening at: 127.0.0.1:4650
{"level":"info","ts":1708362016.8454678,"caller":"server/server.go:49","msg":"Serving gRPC 0.0.0.0:8080"}
{"level":"info","ts":1708362016.84577,"caller":"server/server.go:81","msg":"Serving gRPC-Gateway http://0.0.0.0:8090"}
{"level":"info","ts":1708362017.848982,"caller":"server/service.go:44","msg":"received request","endpoint":"CurrentWeather","lat":51.50732,"lon":-0.1276474}
res.Status: 200 OK
body: {"condition":"Clouds","climate":"moderate"}
PASS
*/
func TestServerValidRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go Serve(ctx)
	//wait for server to start
	<-time.After(time.Second)

	// build request
	req, err := http.NewRequest(http.MethodGet, "http://0.0.0.0:8090/v1/weather", nil)
	assert.NoError(t, err)
	q := req.URL.Query()
	q.Add("lat", "51.5073219")
	q.Add("lon", "-0.1276474")
	req.URL.RawQuery = q.Encode()

	// send request
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	// read results
	resBytes, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	fmt.Printf("res.Status: %v\n", res.Status)
	fmt.Printf("body: %v\n", string(resBytes))

}

/*
API server listening at: 127.0.0.1:21887
{"level":"info","ts":1708362243.804499,"caller":"server/server.go:49","msg":"Serving gRPC 0.0.0.0:8080"}
{"level":"info","ts":1708362243.804974,"caller":"server/server.go:81","msg":"Serving gRPC-Gateway http://0.0.0.0:8090"}
{"level":"info","ts":1708362244.807452,"caller":"server/service.go:44","msg":"received request","endpoint":"CurrentWeather","lat":1111111200,"lon":2222222300}
{"level":"info","ts":1708362244.807709,"caller":"server/service.go:48","msg":"invalid request","validation_error":"Key: 'Lat' Error:Field validation for 'Lat' failed on the 'latitude' tag\nKey: 'Lon' Error:Field validation for 'Lon' failed on the 'longitude' tag"}
res.Status: 400 Bad Request
body: {"code":3,"message":"Key: 'Lat' Error:Field validation for 'Lat' failed on the 'latitude' tag\nKey: 'Lon' Error:Field validation for 'Lon' failed on the 'longitude' tag","details":[]}
PASS
*/
func TestServerInvalidRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go Serve(ctx)
	//wait for server to start
	<-time.After(time.Second)

	// build request
	req, err := http.NewRequest(http.MethodGet, "http://0.0.0.0:8090/v1/weather", nil)
	assert.NoError(t, err)
	q := req.URL.Query()
	q.Add("lat", "1111111111")
	q.Add("lon", "2222222222")
	req.URL.RawQuery = q.Encode()

	// send request
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	// read results
	resBytes, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	fmt.Printf("res.Status: %v\n", res.Status)
	fmt.Printf("body: %v\n", string(resBytes))
}
