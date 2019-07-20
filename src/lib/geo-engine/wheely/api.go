package wheely

import (
	"context"
	"net/http"
)

type API struct {
	client *http.Client
	cfg    Config
}

func (it *API) DoCars(ctx context.Context, lat, lng float64, any interface{}) error {
	if ctx == nil {
		ctx = context.TODO()
	}
	return nil
}

func (it *API) DoPredict(ctx context.Context) error {
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
