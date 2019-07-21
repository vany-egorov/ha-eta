package cache

import (
	"fmt"

	"github.com/vany-egorov/ha-eta/lib/cache/common"
	goCache "github.com/vany-egorov/ha-eta/lib/cache/go-cache"
)

type config struct {
	Kind Kind `yaml:"kind"`

	Common common.Config `yaml:",inline"`
}

type ConfigInternal interface {
	Defaultize()
}

type Config struct {
	config `yaml:"-"`

	internal ConfigInternal `yaml:"-"`
}

func (it *Config) Defaultize() {
	it.Common.Defaultize()

	if intrl := it.internal; intrl != nil {
		intrl.Defaultize()
	}
}

func (it *Config) WithGoCache(fn func(cfg *goCache.Config)) {
	if v, ok := it.internal.(*goCache.Config); ok {
		fn(v)
	}
}

func (it *Config) InitWithKind(kind Kind) error {
	it.Kind = kind

	switch kind {
	case KindGoCache:
		it.internal = new(goCache.Config)

	default:
		return fmt.Errorf(
			"got '%s' kind while initializing cache config: only '%s' is supported",
			kind,
			KindGoCache,
		)
	}

	return nil
}
