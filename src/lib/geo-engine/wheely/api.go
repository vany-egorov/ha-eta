package wheely

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	bufPool "github.com/vany-egorov/ha-eta/lib/buf-pool"
)

type API struct {
	client *http.Client
	cfg    Config
}

func (it *API) CarsLimit() uint64 { return it.cfg.CarsLimit }

func (it *API) DoCars(ctx context.Context, lat, lng float64, limit uint64, any interface{}, events interface{}) error {
	if ctx == nil {
		ctx = context.TODO()
	}

	method := methodCars
	u := url.URL{}

	if limit == 0 {
		limit = it.cfg.CarsLimit
	}

	{ // construct url
		patchURLCars(&u, it.cfg.Url)

		q := u.Query()
		q.Set("lat", strconv.FormatFloat(lat, 'f', -1, 64))
		q.Set("lng", strconv.FormatFloat(lng, 'f', -1, 64))
		q.Set("limit", strconv.FormatUint(limit, 10))
		u.RawQuery = q.Encode()
	}

	respBody := bufPool.NewBuf()
	defer respBody.Release()

	if err := doReqRes(
		ctx, it.client,

		method, &u, nil, respBody,

		nil,
		func(respStatusCode int) bool { return respStatusCode == http.StatusOK },
		tryLogReqFn(events),
		tryLogReqResFn(events),
		tryLogRespFn(events),
	); err != nil {
		return err
	}

	cars := Cars{}
	if err := json.Unmarshal(respBody.Bytes(), &cars); err != nil {
		return errors.Wrap(ErrResponseDecode, err.Error())
	}

	cars.mustTo(any)
	return nil
}

func (it *API) DoPredict(ctx context.Context, lat, lng float64, anySrc, anyDst, events interface{}) error {
	if ctx == nil {
		ctx = context.TODO()
	}

	method := methodPredict
	u := url.URL{}

	{ // construct url
		patchURLPredict(&u, it.cfg.Url)
	}

	points := Points{}
	points.mustFrom(anySrc)
	req := PredictReq{Point{lat, lng}, points}

	respBody := bufPool.NewBuf()
	defer respBody.Release()

	if err := doReqRes(
		ctx, it.client,

		method, &u, &req, respBody,

		nil,
		func(respStatusCode int) bool { return respStatusCode == http.StatusOK },
		tryLogReqFn(events),
		tryLogReqResFn(events),
		tryLogRespFn(events),
	); err != nil {
		return err
	}

	var etas ETAs = nil
	if err := json.Unmarshal(respBody.Bytes(), &etas); err != nil {
		return errors.Wrap(ErrResponseDecode, err.Error())
	}

	etas.mustTo(anyDst)
	return nil
}

func (it *API) Initialize(cfg *Config, fnArgs ...interface{}) {
	if cfg != nil {
		it.cfg = *cfg
	}

	it.cfg.Defaultize()

	for _, ifn := range fnArgs {
		if fn, ok := ifn.(Arg); ok {
			fn(&it.cfg)
		}
	}

	it.client = it.cfg.httpClient()
}
