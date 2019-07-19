package servers

import (
	"fmt"
	"path/filepath"

	"github.com/vany-egorov/ha-eta/lib/helpers"
)

type ConfigLogger interface {
	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
}

type Configer interface {
	ToLog(ConfigLogger)
	Defaultize() error
	GetAddr() string
}

type configUnmarshal struct {
	Kind Kind `yaml:"kind"`
}

type Config struct {
	kind Kind

	inner Configer
}

func (it *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	that := configUnmarshal{}
	if e := unmarshal(&that); e != nil {
		return e
	}

	switch that.Kind {
	case KindINET, KindUnknown:
		it.kind = KindINET
		it.inner = &ConfigINET{}
	case KindUNIX:
		it.kind = KindUNIX
		it.inner = &ConfigUNIX{}
	}

	if e := unmarshal(it.inner); e != nil {
		return e
	}

	return nil
}

func (it *Config) Validate() error {
	if it.kind == KindUNIX {
		that := it.inner.(*ConfigUNIX)

		if v := that.Addr; v != "" {
			dir := filepath.Dir(v)
			if exists, e := helpers.IsPathExists(dir); e != nil {
				return fmt.Errorf("unix socket parent dir '%s' existence check failed: %s", dir, e)
			} else if !exists {
				return fmt.Errorf("unix socket parent dir '%s' does not exist", dir)
			}
		} else {
			return fmt.Errorf("unix socket path is empty")
		}
	} else if it.kind == KindINET {
		that := it.inner.(*ConfigINET)

		if that.TLS.Enable {
			if v := that.TLS.CertFile; v != "" {
				if exists, e := helpers.IsPathExists(v); e != nil {
					return fmt.Errorf("tls cert-file existence check failed: %s", e)
				} else if !exists {
					return fmt.Errorf("tls cert-file '%s' does not exist", v)
				}
			} else {
				return fmt.Errorf("tls cert-file path is not provided")
			}

			if v := that.TLS.KeyFile; v != "" {
				if exists, e := helpers.IsPathExists(v); e != nil {
					return fmt.Errorf("tls key-file existence check failed: %s", e)
				} else if !exists {
					return fmt.Errorf("tls key-file '%s' does not exist", v)
				}
			} else {
				return fmt.Errorf("tls key-file path is not provided")
			}
		}
	}

	return nil
}

func (it *Config) RunLogPrefix() string {
	switch it.kind {
	case KindUnknown:
		return "!!!"
	case KindINET:
		that := it.inner.(*ConfigINET)

		if that.TLS.Enable {
			return "HTTPS"
		}

		return "HTTP"
	case KindUNIX:
		return "UNIX"
	}

	return "???"
}

func (it *Config) ToLog(log ConfigLogger) {
	f := log.Debugf

	f("      kind: %s", it.kind)
	it.inner.ToLog(log)
}
