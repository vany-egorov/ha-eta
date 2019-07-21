package errors

import (
	"fmt"
	"net/http"
)

type Err uint

const (
	Unknown Err = iota

	ETAMinReqParse
	ETAMinGeoEngineCars
	ETAMinGeoEnginePredict
	ETAMinNoETAsFound
)

var (
	errText = map[Err]string{
		ETAMinReqParse:         "api/v1:(eta/min): error parse or query-values missing",
		ETAMinGeoEngineCars:    "api/v1:(eta/min): geo-engine call error",
		ETAMinGeoEnginePredict: "api/v1:(eta/min): geo-engine predict error",
		ETAMinNoETAsFound:      "api/v1:(eta/min): no ETAs found",
	}

	errHTTPStatus = map[Err]int{
		Unknown: http.StatusInternalServerError,

		ETAMinReqParse:         http.StatusBadRequest,
		ETAMinGeoEngineCars:    http.StatusBadGateway,
		ETAMinGeoEnginePredict: http.StatusBadGateway,
		ETAMinNoETAsFound:      http.StatusNotFound,
	}
)

func (it Err) String() string {
	if s, ok := errText[it]; ok {
		return s
	}
	return fmt.Sprintf("missing error-code text for 0x%X error", uint(it))
}
func (it Err) Error() string { return it.String() }

/* impl HTTPStatuser for Err */
func (it Err) HTTPStatus() int {
	if s, ok := errHTTPStatus[it]; ok {
		return s
	}
	return http.StatusInternalServerError
}

/* impl APICoder for Err */
func (it Err) APICode() uint { return uint(it) }
