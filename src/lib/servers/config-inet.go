package servers

import (
	"net"
	"strconv"
)

type ConfigINET struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	TLS  struct {
		Enable   bool   `yaml:"enable"`
		CertFile string `yaml:"cert-file"`
		KeyFile  string `yaml:"key-file"`
	} `yaml:"tls"`
}

func (it *ConfigINET) GetAddr() string { return it.Addr() }

func (it *ConfigINET) Addr() string {
	return net.JoinHostPort(it.Host, strconv.Itoa(it.Port))
}

func (it *ConfigINET) ToLog(log ConfigLogger) {
	f := log.Debugf

	f("      host: %s", it.Host)
	f("      port: %d", it.Port)
	if it.TLS.Enable {
		f("      tls:")
		f("        enable: %t", it.TLS.Enable)
		f("        cert-file: %s", it.TLS.CertFile)
		f("        key-file: %s", it.TLS.KeyFile)
	}
}

func (it *ConfigINET) Defaultize() error {
	return nil
}
