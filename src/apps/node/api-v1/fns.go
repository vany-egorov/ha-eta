package v1

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type httpStatuser interface {
	HTTPStatus() int
}

func httpStatus(any error) int {
	if it, ok := any.(httpStatuser); ok {
		return it.HTTPStatus()
	}

	return http.StatusInternalServerError
}

type apiCoder interface {
	APICode() uint
}

func apiCode(any error) uint {
	if it, ok := any.(apiCoder); ok {
		return it.APICode()
	}

	return 666
}

func Send(rw http.ResponseWriter, err error) {
	resp := Resp{
		Code:    apiCode(errors.Cause(err)),
		Message: err.Error(),
	}

	if body, e := json.Marshal(&resp); e != nil {
		http.Error(rw, e.Error(), httpStatus(e))
	} else {
		http.Error(rw, string(body), httpStatus(errors.Cause(err)))
	}
}
