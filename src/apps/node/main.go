package node

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	cli "gopkg.in/urfave/cli.v1"

	apiV1 "github.com/vany-egorov/ha-eta/apps/node/api-v1/handlers"
	"github.com/vany-egorov/ha-eta/handlers"
	"github.com/vany-egorov/ha-eta/lib/gin-contrib/prefix"
)

type App struct {
}

func (it *App) start() (outErr error) {
	router := it.NewRouter()
	server := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	stop := make(chan struct{})
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
				log.Infof("(SIGINT SIGTERM SIGQUIT) will shutdown")
			}
		case <-stop:
		}

		cancel()
	}()

	wg.Add(1)
	go func() {
		defer func() { <-ctx.Done(); wg.Done() }()

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: config
		defer cancel()

		if err := server.Shutdown(ctxTimeout); err != nil {
			log.Errorf("server shutdown failed: %s", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: config
		defer cancel()

		memstats.StartWithCtx(ctx)
		memstats.DoneWithCtx(ctxTimeout)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		close(stop)
		outErr = err
	}

	<-ctx.Done()
	wg.Wait()

	return outErr
}

func (it *App) initialize(ctx *cli.Context) error {
	return nil
}

func (it *App) main() error {
	if err := it.initialize(cliCtx); err != nil {
		return fmt.Errorf("intialization failed: %s\n", e.Error())
	}

	return it.start()
}

func (it *App) Main(cliCtx *cli.Context) {
	if e := it.main(cliCtx); e != nil {
		fmt.Fprint(os.Stderr, e.Error())
	}
}
