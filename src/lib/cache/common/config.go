package common

import "time"

const (
	DefaultPointsTTL = 10 * time.Second
	DefaultETAsTTL   = 10 * time.Second

	DefaultCleanUpInterval = 10 * time.Minute
)

type Config struct {
	DoPoints bool `yaml:"do-points"`
	DoETAs   bool `yaml:"do-etas"`

	PointsTTL time.Duration `yaml:"points-ttl"`
	ETAsTTL   time.Duration `yaml:"etas-ttl"`

	CleanUpInterval time.Duration `yaml:"clean-up-interval"`
}

func (it *Config) Defaultize() {
	it.DoPoints = true
	it.DoETAs = true

	if it.PointsTTL == 0 {
		it.PointsTTL = DefaultPointsTTL
	}

	if it.ETAsTTL == 0 {
		it.ETAsTTL = DefaultETAsTTL
	}

	if it.CleanUpInterval == 0 {
		it.CleanUpInterval = DefaultCleanUpInterval
	}
}
