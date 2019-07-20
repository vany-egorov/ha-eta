package errors

import (
	"fmt"
	"net/http"
)

type Err uint

const (
	Unknown Err = iota

	ETAMinReqParse
)

var (
	errText = map[Err]string{
		ETAMinReqParse: "api/v1:(eta/min): error parse",
	}

	errHTTPStatus = map[Err]int{
		Unknown: http.StatusInternalServerError,

		ETAMinReqParse: http.StatusBadRequest,
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
