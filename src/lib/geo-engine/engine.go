package geoEngine

import (
	"context"
	"fmt"

	"github.com/vany-egorov/ha-eta/lib/geo-engine/wheely"
)

type Engine interface {
	DoCars(ctx context.Context, lat, lng float64, limit uint64, any, events interface{}) error
	DoPredict(ctx context.Context, lat, lng float64, anySrc, anyDst, events interface{}) error
}

func NewGeoEngine(cfg *Config) (Engine, error) {
	switch cfg.Kind {
	case KindWheely:
		engine := new(wheely.API)
		engine.Initialize(cfg.internal.(*wheely.Config))
		return engine, nil
	}

	return nil, fmt.Errorf(
		"got '%s' kind while constructing geo-engine: only '%s' is supported",
		KindWheely,
	)
}
