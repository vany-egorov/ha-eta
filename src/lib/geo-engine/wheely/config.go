package wheely

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	Url *url.URL

	Timeout                   time.Duration `yaml:"timeout"`
	MaxIdleConnectionsPerHost int           `yaml:"max-idle-connections-per-host"`
	DialTimeout               time.Duration `yaml:"dial-timeout"`
	DialKeepAlive             time.Duration `yaml:"dial-keep-alive"`
	DisableKeepAlives         bool          `yaml:"disable-keep-alives"`
	TLSHandshakeTimeout       time.Duration `yaml:"tls-handshake-timeout"`
	TLSInsecureSkipVerify     bool          `yaml:"tls-insecure-skip-verify"`
}

func (it *Config) Defaultize() {
	if it.Url == nil {
		it.Url = defaultUrl
	}

	if it.Timeout == 0 {
		it.Timeout = defaultTimeout
	}
	if it.MaxIdleConnectionsPerHost == 0 {
		it.MaxIdleConnectionsPerHost = defaultMaxIDLEConnectionsPerHost
	}
	if it.DialTimeout == 0 {
		it.DialTimeout = defaultDialTimeout
	}
	if it.DialKeepAlive == 0 {
		it.DialKeepAlive = defaultDialKeepAlive
	}
	if it.TLSHandshakeTimeout == 0 {
		it.TLSHandshakeTimeout = defaultTLSHandshakeTimeout
	}
}

func (it *Config) httpClient() *http.Client {
	return &http.Client{
		Timeout: it.Timeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: it.MaxIdleConnectionsPerHost,
			Dial: (&net.Dialer{
				Timeout:   it.DialTimeout,
				KeepAlive: it.DialKeepAlive,
			}).Dial,
			TLSHandshakeTimeout: it.TLSHandshakeTimeout,
			DisableKeepAlives:   it.DisableKeepAlives,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: it.TLSInsecureSkipVerify,
			},
		},
	}
}

type Arg func(*Config)

func Url(v *url.URL) Arg {
	return func(cfg *Config) { cfg.Url = v }
}

func Timeout(v time.Duration) Arg {
	return func(cfg *Config) { cfg.Timeout = v }
}

func MaxIdleConnectionsPerHost(v int) Arg {
	return func(cfg *Config) { cfg.MaxIdleConnectionsPerHost = v }
}

func DialTimeout(v time.Duration) Arg {
	return func(cfg *Config) { cfg.DialTimeout = v }
}

func DialKeepAlive(v time.Duration) Arg {
	return func(cfg *Config) { cfg.DialKeepAlive = v }
}

func DisableKeepAlives(v bool) Arg {
	return func(cfg *Config) { cfg.DisableKeepAlives = v }
}

func TLSHandshakeTimeout(v time.Duration) Arg {
	return func(cfg *Config) { cfg.TLSHandshakeTimeout = v }
}

func TLSInsecureSkipVerify(v bool) Arg {
	return func(cfg *Config) { cfg.TLSInsecureSkipVerify = v }
}
