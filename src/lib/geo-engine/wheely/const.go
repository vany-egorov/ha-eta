package wheely

import (
	"net/http"
	"net/url"
	"time"
)

const (
	methodCars = http.MethodGet
	pathCars   = "/cars"

	methodPredict = http.MethodPost
	pathPredict   = "/predict"
)

var (
	defaultUrl    *url.URL = nil
	DefaultUrlRaw string   = "https://dev-api.wheely.com/fake-eta"

	DefaultCarsLimit uint64 = 10

	defaultTimeout                   time.Duration = 1 * time.Second
	defaultMaxIDLEConnectionsPerHost int           = 16
	defaultDialTimeout               time.Duration = 1 * time.Second
	defaultDialKeepAlive             time.Duration = 10 * time.Second
	defaultDisableKeepAlives         bool          = false
	defaultTLSHandshakeTimeout       time.Duration = 1 * time.Second
	defaultTLSInsecureSkipVerify     bool          = true
)

func init() {
	// must-parse!
	if v, err := url.Parse(DefaultUrlRaw); err != nil {
		panic(err)
	} else {
		defaultUrl = v
	}
}
