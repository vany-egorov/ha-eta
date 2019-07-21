package wheely

import "errors"

var (
	ErrRequestEncode  = errors.New("error encode request data")
	ErrRequestCreate  = errors.New("error create http request")
	ErrRequestExecute = errors.New("error execute http request")
	ErrResponseRead   = errors.New("error read data from response body into buffer")
	BadStatusCode     = errors.New("bad http-status code")
	ErrResponseDecode = errors.New("error decode response body")
)
