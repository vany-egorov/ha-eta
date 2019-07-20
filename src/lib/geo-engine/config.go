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

	Internal ConfigInternal `yaml:"-"`
}

func (it *Config) Defaultize() {
	if intrl := it.Internal; intrl != nil {
		intrl.Defaultize()
	}
}

func (it *Config) InitWithKind(kind Kind) error {
	it.Kind = kind

	switch kind {
	case KindWheely:
		it.Internal = new(wheely.Config)

	default:
		return fmt.Errorf(
			"got '%s' kind while initializing geo-engine config, only '%s' is supported",
			KindWheely,
		)
	}

	return nil
}
