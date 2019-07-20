package cache

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
