package node

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/vany-egorov/ha-eta/lib/cache"
	geoEngine "github.com/vany-egorov/ha-eta/lib/geo-engine"
	"github.com/vany-egorov/ha-eta/lib/geo-engine/wheely"
)

type config struct {
	Server struct {
		Host string
		Port int
	}

	Period struct {
		Memstats time.Duration
	}

	Timeout struct {
		WaitTerminate time.Duration
	}

	GeoEngine geoEngine.Config
	Cache     cache.Config
}

func (it *config) serverAddr() string {
	return net.JoinHostPort(it.Server.Host,
		strconv.Itoa(it.Server.Port))
}

func (it *config) parseFromFile() error { return nil }

func (it *config) parse(c *cli.Context) error {
	if err := it.parseFromFile(); err != nil {
		return err
	}

	if c.IsSet("server-port") {
		it.Server.Port = c.Int("server-port")
	}

	if c.IsSet("server-host") {
		it.Server.Host = c.String("server-host")
	}

	if c.IsSet("period-memstats") {
		it.Period.Memstats = c.Duration("period-memstats")
	}

	if c.IsSet("timeout-wait-terminate") {
		it.Timeout.WaitTerminate = c.Duration("timeout-wait-terminate")
	}

	if c.IsSet("geo-engine-kind") {
		raw := c.String("geo-engine-kind")
		kind := geoEngine.NewKindFromString(raw)
		if kind == geoEngine.KindUnknown {
			return fmt.Errorf("error parse geo-engine-kind from: %s", raw)
		}

		if kind != it.GeoEngine.Kind {
			it.GeoEngine.InitWithKind(kind)
		}
	}

	if c.IsSet("wheely-url") {
		raw := c.String("wheely-url")
		u, err := url.Parse(raw)
		if err != nil {
			return fmt.Errorf("error parse wheely-url: %s", err)
		}

		it.GeoEngine.WithWheely(func(cfg *wheely.Config) {
			cfg.Url = u
		})
	}

	if c.IsSet("wheely-cars-limit") {
		v := c.Uint64("wheely-cars-limit")
		it.GeoEngine.WithWheely(func(cfg *wheely.Config) {
			cfg.CarsLimit = v
		})
	}

	return nil
}

func (it *config) validate(actn action) error { return nil }

func (it *config) defaultize() {
	it.Server.Host = defaultServerHost
	it.Server.Port = defaultServerPort
	it.Period.Memstats = defaultPeriodMemstats
	it.Timeout.WaitTerminate = defaultTimeoutWaitTerminate

	it.GeoEngine.InitWithKind(geoEngine.DefaultKind)
	it.GeoEngine.Defaultize()
}

func (it *config) build(c *cli.Context, actn action) error {
	it.defaultize()

	if err := it.parse(c); err != nil {
		return err
	}

	if err := it.validate(actn); err != nil {
		return err
	}

	return nil
}
