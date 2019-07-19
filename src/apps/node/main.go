package node

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	daemon "github.com/sevlyar/go-daemon"
	"github.com/wsxiaoys/terminal/color"
	cli "gopkg.in/urfave/cli.v1"

	apiV1 "github.com/vany-egorov/ha-eta/apps/node/api-v1/handlers"
	appCtx "github.com/vany-egorov/ha-eta/apps/node/ctx"
	"github.com/vany-egorov/ha-eta/apps/node/initializers"
	"github.com/vany-egorov/ha-eta/handlers"
	"github.com/vany-egorov/ha-eta/lib"
	"github.com/vany-egorov/ha-eta/lib/gin-contrib/prefix"
	"github.com/vany-egorov/ha-eta/lib/gin-contrib/seelog"
	"github.com/vany-egorov/ha-eta/lib/memstats"
)

type App struct {
	daemonContext *daemon.Context

	memstats *memstats.Memstats
}

func (it *App) SetVersion(v string)     { appCtx.Ctx().SetVersion(v) }
func (it *App) SetDoneAt(v *lib.DoneAt) { appCtx.Ctx().SetDoneAt(v) }

func (it *App) ShowVersion(ctx *cli.Context) {
	fmt.Printf("%s\n", initializers.Logo)
	fmt.Printf(color.Sprintf("@m%s\n", initializers.ServiceName))
	fmt.Printf(color.Sprintf("@m%s\n\n", initializers.Company))
	fmt.Printf(color.Sprintf("version: @y%s\n\n", appCtx.Ctx().Version()))
	appCtx.Ctx().DoneAt().Print()
}
func (it *App) ShowVersionAndExit(ctx *cli.Context) {
	it.ShowVersion(ctx)
	os.Exit(0)
}

func (it *App) Daemonize() error {
	it.daemonContext = &daemon.Context{
		Umask: appCtx.Ctx().Config().Daemon.Umask,
		Args:  os.Args,
	}

	if appCtx.Ctx().Config().Daemon.Pidfile != "" {
		it.daemonContext.PidFileName = appCtx.Ctx().Config().Daemon.Pidfile
		it.daemonContext.PidFilePerm = appCtx.Ctx().Config().Daemon.PidfileMode
	}

	if appCtx.Ctx().Config().Daemon.WorkDir != "" {
		it.daemonContext.WorkDir = appCtx.Ctx().Config().Daemon.WorkDir
	}

	child, e := it.daemonContext.Reborn()
	if e != nil {
		return fmt.Errorf("daemon reborn failed: %s", e.Error())
	}

	if child != nil { // parent exit
		os.Exit(0)
	}

	return nil
}

func (it *App) NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	ginHTTPLogger := seelog.New(appCtx.HTTPLogger())

	r := gin.New()
	r.Use(
		gin.Recovery(),
		prefix.New(),
	)

	if cfg := appCtx.Ctx().Config().CORS; cfg.Enable {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowMethods = cfg.AllowMethods
		corsConfig.AllowHeaders = cfg.AllowHeaders
		if cfg.AllowAllOrigins {
			corsConfig.AllowAllOrigins = true
		} else {
			corsConfig.AllowOriginFunc = func(origin string) bool { return true }
		}
		corsConfig.AllowCredentials = cfg.AllowCredentials

		r.Use(
			// log only OPTION requests
			seelog.NewCommmon(appCtx.HTTPLogger(), func(c *gin.Context) bool { return c.Request.Method == "OPTIONS" }),
			cors.New(corsConfig),
		)
	}

	r.Use(func(c *gin.Context) { c.Set("app-cxt", appCtx.Ctx()); c.Next() })

	r.HandleMethodNotAllowed = true
	r.NoRoute(handlers.NoRoute, ginHTTPLogger)
	r.NoMethod(handlers.NoMethod, ginHTTPLogger)

	{
		α := r.Group("/api/v1")
		α.Use(ginHTTPLogger)

		α.GET("/eta/min", apiV1.ETAMin)
	}

	return r
}

func (it *App) Start() {
	log.Debugf("=> started - app")
	log.Infof("application '%s' with pid=%d started in %s environment",
		initializers.ServiceName, os.Getpid(), appCtx.Ctx().Config().Environment)
	if appCtx.Ctx().Config().NoConfigFileWasParsed() {
		log.Infof("=> no configuration file was parsed - using default configuration")
	} else {
		log.Infof(`=> using configuration from "%s"`, appCtx.Ctx().Config().PathConfig())
	}
	if appCtx.Ctx().Config().Log.IsOutputToFile() {
		log.Infof(`=> log directory "%s"`, appCtx.Ctx().Config().Log.Path)
	} else {
		log.Infof("=> no output to file. console output only")
	}

	if appCtx.Ctx().Config().Daemonize && appCtx.Ctx().Config().Daemon.Pidfile != "" {
		log.Infof(`=> writing pidfile to "%s" as %s`,
			appCtx.Ctx().Config().Daemon.Pidfile, appCtx.Ctx().Config().Daemon.PidfileMode)
	}

	appCtx.Ctx().Config().ToLog(log.Current)

	it.memstats = memstats.NewMemstats(appCtx.Ctx()).GoStart()

	{ // http
		router := it.NewRouter()
		appCtx.Ctx().Config().Servers.RunGin(router, log.Current)
	}
}

func (it *App) init(ctx *cli.Context) error {
	flags := initializers.NewFlags(ctx)
	if flags.ShowVersionAndExit() {
		it.ShowVersionAndExit(ctx)
	}
	appCtx.Ctx().SetFlags(flags)

	c, e := initializers.NewConfig(flags)
	if c != nil && c.PrintConfig() {
		if e != nil {
			fmt.Fprintf(os.Stderr, "initialization failed: %s\n", color.Sprintf("@r%s", e.Error()))
		}
		if c.NoConfigFileWasParsed() {
			fmt.Printf("=> no configuration file was parsed (tried \"%s\") - using default configuration\n", c.PathConfig())
		} else {
			fmt.Printf("=> using configuration from \"%s\"\n", c.PathConfig())
		}
		c.ToLog(nil)
		os.Exit(0)
	}

	if e != nil {
		return e
	} else {
		appCtx.Ctx().SetConfig(c)
	}

	if appCtx.Ctx().Config().Daemonize {
		if e := it.Daemonize(); e != nil {
			return e
		}
	}

	if e := appCtx.LoadLoggers(); e != nil {
		return e
	}

	log.Debugf("=> successfully initialized")

	return nil
}

func (it *App) signalHadler() {
	log.Debugf("=> started - signal handler")
	log.Debugf("=> SIGINT|SIGTERM|SIGQUIT - stop")
	log.Debugf("=> SIGUSR1                - postrotate | reload loggers config")
	log.Debugf("=> SIGHUP                 - reload config")
	for {
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan,
			syscall.SIGINT, os.Interrupt, // CTRL-C
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGUSR1, // postrotate
			syscall.SIGHUP,  // reload
		)

		for sig := range signalChan {
			func() {
				defer log.Flush()
				switch sig {
				case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
					log.Infof("[SIGINT, SIGTERM, SIGQUIT] will shutdown;")

					it.memstats.Stop()

					allDone := make(chan struct{}, 1)
					go func() {
						defer func() { allDone <- struct{}{} }()
						it.memstats.Done()
					}()

					select {
					case <-allDone:
					case <-time.After(2 * time.Second):
					}

					if it.daemonContext != nil {
						it.daemonContext.Release()
						it.daemonContext = nil
					}

					if appCtx.Ctx().Config().Daemonize && appCtx.Ctx().Config().Daemon.Pidfile != "" {
						os.Remove(appCtx.Ctx().Config().Daemon.Pidfile)
						log.Infof("[SIGINT, SIGTERM, SIGQUIT] => OK removing pidfile %s", appCtx.Ctx().Config().Daemon.Pidfile)
					}
					log.Infof("application with pid=%d stopped;\n", os.Getpid())
					log.Flush()

					os.Exit(0)

				case syscall.SIGUSR1:
					defer appCtx.Ctx().DoneAt().UpdatePostrotatedAt()
					log.Infof("[SIGUSR1] => postrotate logs reloading")
					log.Infof("[SIGUSR1] => reloading configuration from config file")
					if c, e := initializers.NewConfig(appCtx.Ctx().Flags()); e != nil {
						log.Errorf("[SIGUSR1] => reloading config failed: %s", e.Error())
					} else {
						log.Infof("[SIGUSR1] => replacing old log configuration")
						appCtx.Ctx().Config().SetLog(c.GetLog())
						log.Infof("[SIGUSR1] => reloading loggers")
						if e := appCtx.LoadLoggers(); e != nil {
							log.Errorf("[SIGUSR1] => reloading config failed: %s", e.Error())
						} else {
							log.Infof("[SIGUSR1] => loggers and log settings reloaded successfully. new loggers configuration is:")
							appCtx.Ctx().Config().GetLog().ToLog(log.Current)
						}
					}

				case syscall.SIGHUP:
					defer appCtx.Ctx().DoneAt().UpdateReloadedAt()
					log.Infof("[SIGHUP] => reloading")
					log.Infof("[SIGHUP] => reloading configuration file")
					if c, e := initializers.NewConfig(appCtx.Ctx().Flags()); e != nil {
						log.Errorf("[SIGHUP] => reloading config failed: %s", e.Error())
					} else {
						log.Infof("[SIGHUP] => replacing old configuration")
						appCtx.Ctx().SetConfig(c)
						log.Infof("[SIGHUP] => reloading loggers")
						if e := appCtx.LoadLoggers(); e != nil {
							log.Errorf("[SIGHUP] => reloading config failed: %s", e.Error())
						} else {
							log.Infof("[SIGHUP] => configuration reloaded successfully. new configuration is:")
							appCtx.Ctx().Config().ToLog(log.Current)
						}

						it.memstats.Stop()

						it.memstats.Done()

						it.memstats = memstats.NewMemstats(appCtx.Ctx()).GoStart()
					}
				}
			}()
		}
	}
}

func (it *App) Main(ctx *cli.Context, action Action) error {
	if e := it.init(ctx); e != nil {
		fmt.Fprintf(os.Stderr, "bootstrap failed: %s\n", color.Sprintf("@r%s", e.Error()))
		os.Exit(1)
	}

	var e error = nil
	switch action {
	case ActionMain:
		go it.Start()
		it.signalHadler()
	default:
		fmt.Fprintf(os.Stderr, "action \"%s\" is not implemented\n", action)
	}

	if e != nil {
		fmt.Fprintf(os.Stderr, "%s failed: %s\n", action, color.Sprintf("@r%s", e.Error()))
		os.Exit(1)
	}

	return nil
}

func NewApp() *App { return new(App) }
