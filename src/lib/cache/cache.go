package cache

import (
	"fmt"

	goCache "github.com/vany-egorov/ha-eta/lib/cache/go-cache"
	"github.com/vany-egorov/ha-eta/models"
)

type Cache interface {
	GetPoints(point models.Point, limit uint64, points *models.Points) bool
	SetPoints(point models.Point, limit uint64, points models.Points)

	GetETAs(point models.Point, all models.Points, hits, miss *models.Points, etas *models.ETAs)
	SetETAs(point models.Point, points models.Points, etas models.ETAs)

	Flush()
}

func NewCache(cfg *Config) (Cache, error) {
	switch cfg.Kind {
	case KindGoCache:
		cache := new(goCache.Cache)
		cache.Initialize(cfg.internal.(*goCache.Config), cfg.Common)
		return cache, nil
	}

	return nil, fmt.Errorf(
		"got '%s' kind while constructing cache: only '%s' is supported",
		cfg.Kind,
		KindGoCache,
	)
}
