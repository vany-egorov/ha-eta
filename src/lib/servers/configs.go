package servers

import (
	"os"
	"time"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

type Configs []*Config

func (it *Configs) PushINETIfNotExists(host string, port int) {
	for _, c := range *it {
		if c.kind == KindINET {
			return
		}
	}

	that := &Config{
		kind: KindINET,
		inner: &ConfigINET{
			Host: host,
			Port: port,
		},
	}

	*it = append(*it, that)
}

func (it *Configs) PushUnixIfNotExists(addr string) {
	for _, c := range *it {
		if c.kind == KindUNIX {
			return
		}
	}

	that := &Config{
		kind: KindUNIX,
		inner: &ConfigUNIX{
			Addr: addr,
		},
	}

	*it = append(*it, that)
}

func (it Configs) Defaultize(inetHost string, inetPort int, unixAddr string) error {
	doneUNIX := false
	doneINET := false

	for _, c := range it {
		if e := c.inner.Defaultize(); e != nil {
			return e
		}

		if c.kind == KindUNIX && !doneUNIX {
			that := c.inner.(*ConfigUNIX)
			if that.Addr == "" {
				that.Addr = unixAddr
			}

			doneUNIX = true
		}

		if c.kind == KindINET && !doneINET {
			that := c.inner.(*ConfigINET)
			if that.Host == "" {
				that.Host = inetHost
			}
			if that.Port == 0 {
				that.Port = inetPort
			}

			doneINET = true
		}
	}

	return nil
}

func (it Configs) Validate() error {
	for _, c := range it {
		if e := c.Validate(); e != nil {
			return e
		}
	}

	return nil
}

func (it Configs) ToLog(log ConfigLogger) {
	f := log.Debugf

	for _, c := range it {
		f("    -")
		c.ToLog(log)
	}
}

func (it Configs) RunGin(router *gin.Engine, log seelog.LoggerInterface) {
	for _, c := range it {
		go func(c *Config) {
			log.Infof("%s server starting on %s", c.RunLogPrefix(), c.inner.GetAddr())
			log.Flush()

			switch c.kind {
			case KindINET:
				that := c.inner.(*ConfigINET)
				addr := that.GetAddr()

				if that.TLS.Enable {
					if e := router.RunTLS(addr, that.TLS.CertFile, that.TLS.KeyFile); e != nil {
						log.Criticalf("starting https (%s) server failed: %s", addr, e.Error())
						os.Exit(1)
					}
				} else {
					if e := router.Run(addr); e != nil {
						log.Criticalf("starting http (%s) server failed: %s", addr, e.Error())
						os.Exit(1)
					}
				}

			case KindUNIX:
				that := c.inner.(*ConfigUNIX)
				addr := that.GetAddr()
				mode := that.SocketFileMode

				go func() {
					<-time.After(1 * time.Second)
					if e := os.Chmod(addr, mode); e != nil {
						log.Criticalf("changing permissions to unix socket %s failed: %s", addr, e.Error())
						os.Exit(1)
					} else {
						log.Infof(`{"status": "chmod OK", "perms": "%03o | %s", "addr": "%s", "cmd": "chmod %o %s"}`,
							mode.Perm(), mode, addr, mode.Perm(), addr)
					}
				}()

				if e := router.RunUnix(addr); e != nil {
					log.Criticalf("starting unix (%s) server failed: %s", addr, e.Error())
					os.Exit(1)
				}
			}
		}(c)
	}
}
