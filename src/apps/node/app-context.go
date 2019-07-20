package node

import (
	"sync"

	cache "github.com/vany-egorov/ha-eta/lib/cache"
	geoEngine "github.com/vany-egorov/ha-eta/lib/geo-engine"
)

type Context struct {
	cfgLock sync.RWMutex
	config  *config

	geoEngineLock sync.RWMutex
	geo           geoEngine.Engine

	cacheLock sync.RWMutex
	cch       cache.Cache
}

func (it *Context) setCfg(v *config) *config {
	it.cfgLock.Lock()
	defer it.cfgLock.Unlock()
	old := it.config
	it.config = v
	return old
}

func (it *Context) cfg() *config {
	it.cfgLock.RLock()
	defer it.cfgLock.RUnlock()
	return it.config
}

func (it *Context) setGeoEngine(v geoEngine.Engine) geoEngine.Engine {
	it.geoEngineLock.Lock()
	defer it.geoEngineLock.Unlock()
	old := it.geo
	it.geo = v
	return old
}

func (it *Context) geoEngine(v geoEngine.Engine) geoEngine.Engine {
	it.geoEngineLock.RLock()
	defer it.geoEngineLock.RUnlock()
	return it.geo
}

func (it *Context) setCache(v cache.Cache) cache.Cache {
	it.cacheLock.Lock()
	defer it.cacheLock.Unlock()
	old := it.cch
	it.cch = v
	return old
}

func (it *Context) cache(v cache.Cache) cache.Cache {
	it.cacheLock.RLock()
	defer it.cacheLock.RUnlock()
	return it.cch
}
