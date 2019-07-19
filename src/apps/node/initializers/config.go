package initializers

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/vany-egorov/ha-eta/lib"
	"github.com/vany-egorov/ha-eta/lib/environment"
	"github.com/vany-egorov/ha-eta/lib/helpers"
	"github.com/vany-egorov/ha-eta/lib/servers"
)

type Config struct {
	flags                 *Flags
	noConfigFileWasParsed bool
	printConfig           bool
	vv                    bool
	vvv                   bool

	Environment environment.Environment `yaml:"environment"`

	Daemonize bool `yaml:"daemonize"`
	Daemon    struct {
		Pidfile     string      `yaml:"pidfile"`
		PidfileMode os.FileMode `yaml:"pidfile-mode"`
		WorkDir     string      `yaml:"workdir"`
		Umask       int         `yaml:"umask"`
	} `yaml:"daemon"`

	Servers servers.Configs `yaml:"servers"`

	CORS struct {
		Enable           bool     `yaml:"enable"`
		AllowMethods     []string `yaml:"allow-methods"`
		AllowHeaders     []string `yaml:"allow-headers"`
		AllowCredentials bool     `yaml:"allow-credentials"`
		AllowAllOrigins  bool     `yaml:"allow-all-origins"`
	} `yaml:"cors"`

	Period struct {
		Memstats time.Duration `yaml:"memstats"`
	} `yaml:"period"`

	Log *lib.ConfigLog `yaml:"log"`
}

func (it *Config) PathLog() string                         { return it.flags.Log() }
func (it *Config) PathLogIsSet() bool                      { return it.flags.ctx.IsSet("log") }
func (it *Config) PathConfig() string                      { return it.flags.Config() }
func (it *Config) NoConfigFileWasParsed() bool             { return it.noConfigFileWasParsed }
func (it *Config) PrintConfig() bool                       { return it.printConfig }
func (it *Config) GetEnvironment() environment.Environment { return it.Environment }
func (it *Config) GetLog() *lib.ConfigLog                  { return it.Log }
func (it *Config) SetLog(v *lib.ConfigLog) *Config         { it.Log = v; return it }

func NewConfig(f *Flags) (*Config, error) {
	it := new(Config)
	it.flags = f

	if it.flags.ctx.IsSet("config") {
		if _, e := os.Stat(it.flags.Config()); e != nil {
			return it, fmt.Errorf("error reading main config file: %s", e.Error())
		}
	}

	if _, e := os.Stat(it.flags.Config()); !os.IsNotExist(e) {
		if in, e := ioutil.ReadFile(it.flags.Config()); e == nil {
			if e := yaml.Unmarshal(in, it); e != nil {
				return nil, fmt.Errorf("yaml.Unmarshal for application config failed: %s", e.Error())
			}
		} else {
			return nil, fmt.Errorf("ioutil.ReadFile for application config file failed: %s", e.Error())
		}
	} else {
		it.noConfigFileWasParsed = true
	}

	if e := it.defaultize(); e != nil {
		return it, e
	}
	it.pathAbsolutize()
	it.pathEnsure()
	it.postProcess()

	if e := it.validate(); e != nil {
		return it, fmt.Errorf("config validation failed: %s", e.Error())
	}

	return it, nil
}

func (it *Config) validate() error {
	if p := it.Log.Path; p != "" {
		if exists, e := helpers.IsPathExists(p); e != nil {
			return fmt.Errorf("Log.Path existence check failed: %s", e)
		} else if !exists {
			return fmt.Errorf("Log.Path: '%s' does not exist", p)
		}
	}

	return nil
}

func (it *Config) defaultize() error {
	if it.Environment.IsUnknown() || it.flags.ctx.IsSet("environment") {
		it.Environment = it.flags.Environment
	}

	{ // daemon
		if it.Daemon.PidfileMode == 0 {
			it.Daemon.PidfileMode = DefaultDaemonPidfileMode
		}

		if it.Daemon.WorkDir == "" {
			it.Daemon.WorkDir = DefaultDaemonWorkdir
		}
	}

	{ // servers
		if e := it.Servers.Defaultize(DefaultServerINETHost, DefaultServerINETPort, DefaultServerUNIXAddr); e != nil {
			return e
		}

		if len(it.Servers) == 0 {
			it.Servers.PushINETIfNotExists(DefaultServerINETHost, DefaultServerINETPort)
		}
	}

	{ // cors
		if it.CORS.AllowMethods == nil {
			it.CORS.AllowMethods = DefaultCORSAllowMethods
		}
		if it.CORS.AllowHeaders == nil {
			it.CORS.AllowHeaders = DefaultCORSAllowHeaders
		}
	}

	{ // period
		if it.Period.Memstats == 0 {
			it.Period.Memstats = DefaultPeriodMemstats
		}
	}

	{ // log
		if it.Log == nil {
			it.Log = lib.NewConfigLog()
		}
		it.Log.Init()

		if it.Log.Path == "" {
			it.Log.Path = DefaultPathLog
		}

		if it.flags.ctx.IsSet("log") || os.Getenv(EnvLog) != "" {
			it.Log.Path = it.flags.Log()
		}

		if !it.Log.Has("app") {
			it.Log.Push("app", "[---]", "info", "critical")
		}

		if !it.Log.Has("http") {
			it.Log.Push("http", "[WS]", "info", "critical")
		}

		if !it.Log.Has("memstats") {
			it.Log.Push("memstats", "[MEMSTATS]", "info", "critical")
		}
	}

	return nil
}

func (it *Config) pathAbsolutize() error { return nil }
func (it *Config) pathEnsure() error     { return nil }
func (it *Config) postProcess() error {
	if it.flags.ctx.IsSet("print-config") && it.flags.ctx.Bool("print-config") {
		it.printConfig = true
	}
	if it.flags.ctx.IsSet("foreground") {
		it.Daemonize = false
	}
	if it.flags.ctx.IsSet("vv") && it.flags.ctx.Bool("vv") {
		it.Log.VV()
	}
	if it.flags.ctx.IsSet("vvv") {
		it.Log.VVV()
	}

	if it.Daemonize {
		it.Log.DisableStdout()
	}

	return nil
}

type ConfigLogger interface {
	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
}

type ConfigLogStdout struct{}

func (it *ConfigLogStdout) Debugf(format string, params ...interface{}) {
	fmt.Printf(format, params...)
	fmt.Println("")
}

func (it *ConfigLogStdout) Infof(format string, params ...interface{}) {
	fmt.Printf(format, params...)
	fmt.Println("")
}

func (it *Config) ToLog(log ConfigLogger) {
	if log == nil {
		log = ConfigLogger(new(ConfigLogStdout))
	}

	f := log.Debugf

	f("config:")
	f("  environment: %s", it.Environment)
	f("  daemonize: %t", it.Daemonize)
	f("  daemon:")
	f("    pidfile: %s", it.Daemon.Pidfile)
	f("    pidfile-mode: %s", it.Daemon.PidfileMode)
	f("    workdir: %s", it.Daemon.WorkDir)
	f("    umask: %03o", it.Daemon.Umask)
	f("  servers:")
	it.Servers.ToLog(log)
	f("  cors:")
	f("    enable: %t", it.CORS.Enable)
	f("    allow-methods: %v", it.CORS.AllowMethods)
	f("    allow-headers: %v", it.CORS.AllowHeaders)
	f("    allow-credentials: %t", it.CORS.AllowCredentials)
	f("    allow-all-origins: %t", it.CORS.AllowAllOrigins)
	f("  period:")
	f("    memstats: %s", it.Period.Memstats)

	it.Log.ToLog(log)
}
