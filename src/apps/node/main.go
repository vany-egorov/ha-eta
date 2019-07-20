package node

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	cli "gopkg.in/urfave/cli.v1"

	"github.com/vany-egorov/ha-eta/lib/cache"
	geoEngine "github.com/vany-egorov/ha-eta/lib/geo-engine"
	"github.com/vany-egorov/ha-eta/lib/log"
	"github.com/vany-egorov/ha-eta/lib/memstats"
)

type App struct {
	ctx Context
}

func (it *App) start() (outErr error) {
	gin.SetMode(gin.ReleaseMode)
	router := it.NewRouter()

	w8Terminate := it.ctx.cfg().Timeout.WaitTerminate

	serverAddr := it.ctx.cfg().serverAddr()
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	serverErrChan := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan,
			syscall.SIGINT, os.Interrupt, // CTRL-C
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		select {
		case _, ok := <-signalChan:
			if ok {
				log.Log(logger(), log.Info, "(SIGINT SIGTERM SIGQUIT) will shutdown")
			}
		case <-serverErrChan:
		}

		cancel()
	}()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		<-ctx.Done()
		defer wg.Done()

		ctxTimeout, cancel := context.WithTimeout(context.Background(), w8Terminate)
		defer cancel()

		if err := server.Shutdown(ctxTimeout); err != nil {
			log.Log(logger(), log.Info, fmt.Sprintf("server shutdown failed: %s", err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		worker := memstats.Worker{}
		worker.Initialize(
			memstats.FnLog(log.LogFnWithLogger(logger())),
			memstats.Period(it.ctx.cfg().Period.Memstats),
		)
		worker.StartWithCtx(ctx)

		ctxTimeout, cancel := context.WithTimeout(context.Background(), w8Terminate)
		defer cancel()
		worker.DoneWithContext(ctxTimeout)
	}()

	log.Log(logger(), log.Info,
		fmt.Sprintf("application with (:pid %d) started", os.Getpid()))

	log.Log(logger(), log.Info,
		fmt.Sprintf("listening on %s", serverAddr))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		close(serverErrChan)
		outErr = err
	}

	wg.Wait()

	return outErr
}

// apply flags, reload config, reload loggers
func (it *App) initialize(cliCtx *cli.Context, actn action) error {
	cfg := new(config)

	if e := cfg.build(cliCtx, actn); e != nil {
		return e
	}

	if e := initLogger(cfg); e != nil {
		return e
	}

	it.ctx.setCfg(cfg)

	if v, e := geoEngine.NewGeoEngine(&cfg.GeoEngine); e != nil {
		it.ctx.setGeoEngine(v)
	}

	if v, e := cache.NewCache(&cfg.Cache); e != nil {
		it.ctx.setCache(v)
	}

	return nil
}

func (it *App) main(cliCtx *cli.Context, actn action) error {
	if err := it.initialize(cliCtx, actn); err != nil {
		return fmt.Errorf("intialization failed: %s\n", err.Error())
	}

	return it.start()
}

func (it *App) Main(cliCtx *cli.Context, actn action) error {
	return it.main(cliCtx, actn)
}
