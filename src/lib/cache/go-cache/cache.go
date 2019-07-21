package goCache

import (
	"sync"

	cache "github.com/patrickmn/go-cache"

	bufPool "github.com/vany-egorov/ha-eta/lib/buf-pool"
	"github.com/vany-egorov/ha-eta/lib/cache/common"
	"github.com/vany-egorov/ha-eta/models"
)

type Cache struct {
	c *cache.Cache

	pointsLock sync.RWMutex

	cfgc common.Config
	cfg  Config
}

func (it *Cache) GetPoints(point models.Point, limit uint64, points *models.Points) bool {
	if !it.cfgc.DoPoints {
		return false
	}

	buf := bufPool.NewBuf()
	defer buf.Release()

	keyPoints(&buf.Buffer, point, limit)
	key := buf.String()

	{ // guard
		it.pointsLock.Lock()
		defer it.pointsLock.Unlock()

		iv, ok := it.c.Get(key)
		if !ok {
			return false
		}

		v, ok := iv.(models.Points)
		if !ok {
			it.c.Delete(key)
			return false
		}

		*points = v
	}

	return true
}

func (it *Cache) SetPoints(point models.Point, limit uint64, points models.Points) {
	if !it.cfgc.DoPoints {
		return
	}

	buf := bufPool.NewBuf()
	defer buf.Release()

	it.pointsLock.RLock()
	defer it.pointsLock.RUnlock()

	keyPoints(&buf.Buffer, point, limit)
	key := buf.String()

	it.c.Set(key, points, it.cfgc.PointsTTL)
}

func (it *Cache) GetETAs(point models.Point, all models.Points, hits, miss *models.Points, etas *models.ETAs) {
	if !it.cfgc.DoETAs {
		*miss = all
		return
	}

	buf := bufPool.NewBuf()
	defer buf.Release()

	for _, pnt := range all {
		keyETAs(&buf.Buffer, point, pnt)
		key := buf.String()
		buf.Reset()

		if iv, ok := it.c.Get(key); ok {
			*hits = append(*hits, pnt)
			*etas = append(*etas, iv.(models.ETA))
		} else {
			*miss = append(*miss, pnt)
		}
	}
}

// TODO: TTL
func (it *Cache) SetETAs(point models.Point, points models.Points, etas models.ETAs) {
	if !it.cfgc.DoETAs {
		return
	}

	buf := bufPool.NewBuf()
	defer buf.Release()

	models.FnZipForEarh(points, etas, func(pnt models.Point, et models.ETA) {
		defer buf.Reset()

		keyETAs(&buf.Buffer, point, pnt)
		key := buf.String()

		it.c.Set(key, et, it.cfgc.ETAsTTL)
	})
}

func (it *Cache) Flush() { it.c.Flush() }

func (it *Cache) Initialize(cfg *Config, cfgCommon common.Config) {
	cfg.Defaultize()

	it.cfgc = cfgCommon
	it.cfg = *cfg

	it.c = cache.New(cache.NoExpiration, it.cfgc.CleanUpInterval)
}
