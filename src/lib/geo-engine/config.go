package geoEngine

import (
	"fmt"

	"github.com/vany-egorov/ha-eta/lib/geo-engine/wheely"
)

type config struct {
	Kind Kind `yaml:"kind"`
}

type ConfigInternal interface {
	Defaultize()
}

type Config struct {
	config `yaml:"-"`

	internal ConfigInternal `yaml:"-"`
}

func (it *Config) Defaultize() {
	if intrl := it.internal; intrl != nil {
		intrl.Defaultize()
	}
}

func (it *Config) WithWheely(fn func(cfg *wheely.Config)) {
	if v, ok := it.internal.(*wheely.Config); ok {
		fn(v)
	}
}

func (it *Config) InitWithKind(kind Kind) error {
	it.Kind = kind

	switch kind {
	case KindWheely:
		it.internal = new(wheely.Config)

	default:
		return fmt.Errorf(
			"got '%s' kind while initializing geo-engine config: only '%s' is supported",
			KindWheely,
		)
	}

	return nil
}
