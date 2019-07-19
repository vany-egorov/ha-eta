package lib

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/cihub/seelog"

	"github.com/vany-egorov/ha-eta/lib/environment"
)

// TODO: different logger-kind for different environment:
//   - dev       - sync;
//   - prod/test - async;
const DefaultXMLTemplate string = `
<seelog
	type="asynctimer" asyncinterval="50000"
	minlevel="{{.Level.Min}}" maxlevel="{{.Level.Max}}"
>
	<outputs>
		{{if .DisableStdout}}{{else}}
			<filter levels="trace">
				<console formatid="default"/>
			</filter>
			<filter levels="debug">
				<console formatid="default"/>
			</filter>
			<filter levels="info">
				<console formatid="default"/>
			</filter>
			<filter levels="warn">
				<console formatid="default"/>
			</filter>
			<filter levels="error">
				<console formatid="default"/>
			</filter>
			<filter levels="critical">
				<console formatid="default"/>
			</filter>
		{{end}}

		{{if .DisableStdout}}{{else}}
			{{if .Path.Error}}
				<filter levels="warn,error,critical">
					<file path="{{.Path.Error}}" formatid="formatFile"/>
				</filter>
			{{end}}
			{{if .Path.Access}}
				<file path="{{.Path.Access}}" formatid="formatFile"/>
			{{end}}
		{{end}}
	</outputs>
	<formats>
		<format id="default" format="[%Date(2/Jan/2006 15:04:05.000)] {{.Prefix}} [%l] %Msg%n"/>
		<format id="formatFile" format="[%Date(2/Jan/2006 15:04:05.000)] {{.Prefix}} [%l] %Msg%n"/>
	</formats>
</seelog>`

type configLogger struct {
	DisableStdout bool   `yaml:"disable-stdout"`
	DisableFile   bool   `yaml:"disable-file"`
	Prefix        string `yaml:"prefix"`
	Path          struct {
		Access string `yaml:"access"`
		Error  string `yaml:"error"`
	} `yaml:"path"`
	Files struct {
		Access string `yaml:"access"`
		Error  string `yaml:"error"`
	} `yaml:"files"`
	Level struct {
		Min string `yaml:"min"`
		Max string `yaml:"max"`
	} `yaml:"level"`
}

func (it *configLogger) IsOutputToFile() bool {
	if it.DisableFile {
		return false
	}

	is1 := it.Files.Access == ""
	is2 := it.Files.Error == ""
	is3 := it.Path.Access == ""
	is4 := it.Path.Error == ""
	if is1 && is2 && is3 && is4 {
		return false
	}
	return true
}

type ConfigLog struct {
	Path    string `yaml:"path"`
	Loggers map[string]*configLogger
}

func (it *ConfigLog) DisableStdout() {
	for _, logger := range it.Loggers {
		logger.DisableStdout = true
	}
}
func (it *ConfigLog) VV() {
	for _, logger := range it.Loggers {
		logger.Level.Min = "debug"
	}
}
func (it *ConfigLog) VVV() {
	for _, logger := range it.Loggers {
		logger.Level.Min = "trace"
	}
}

func (it *ConfigLog) IsOutputToFile() bool {
	for _, it := range it.Loggers {
		if it.IsOutputToFile() {
			return true
		}
	}
	return false
}

type ConfigLogger interface {
	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
}

func (it *ConfigLog) ToLog(log ConfigLogger) {
	f := log.Debugf

	f("  log:")
	f("    path = '%s'", it.Path)
	for loggerName, _ := range it.Loggers {
		loggerConfig := it.Loggers[loggerName]
		f("    %s:", loggerName)
		f("      disable-stdout: %t", loggerConfig.DisableStdout)
		f("      disable-file: %t", loggerConfig.DisableFile)
		f("      prefix: \"%s\"", loggerConfig.Prefix)
		is1 := loggerConfig.Files.Access == ""
		is2 := loggerConfig.Files.Error == ""
		is3 := loggerConfig.Path.Access == ""
		is4 := loggerConfig.Path.Error == ""
		if loggerConfig.DisableFile || is1 && is2 && is3 && is4 {
			f("      no output to file. console output only")
		} else {
			f("      files:")
			f("        access: \"%s\"", loggerConfig.Files.Access)
			f("        error: \"%s\"", loggerConfig.Files.Error)
			f("      path:")
			f("        access: \"%s\"", loggerConfig.Path.Access)
			f("        error: \"%s\"", loggerConfig.Path.Error)
		}
		f("      level:")
		f("        min: \"%s\"", loggerConfig.Level.Min)
		f("        max: \"%s\"", loggerConfig.Level.Max)
	}
}

func (it *ConfigLog) Defaultize(isPathDefaultIsSet bool, pathDefault string, env environment.Environment) {
	if isPathDefaultIsSet {
		it.Path = pathDefault
	}

	// вне зависимости от того, выставлен ли флаг или нет - на prod окружении
	// путь должен быть
	if env.IsProd() && it.Path == "" {
		it.Path = pathDefault
	}

	if it.Path != "" {
		if v, e := filepath.Abs(it.Path); e == nil {
			it.Path = v
		} else {
			it.Path = ""
		}
	}

	for _, loggerConfig := range it.Loggers {
		if loggerConfig.Level.Min == "" {
			loggerConfig.Level.Min = "info"
		}

		if loggerConfig.Level.Max == "" {
			loggerConfig.Level.Max = "critical"
		}

		if loggerConfig.DisableFile {
			continue
		}

		if env.IsProd() { // на production логи должны быть всегда
			if loggerConfig.Files.Access == "" {
				loggerConfig.Files.Access = "access.log"
			}
			if loggerConfig.Files.Error == "" {
				loggerConfig.Files.Error = "error.log"
			}
			if loggerConfig.Path.Access == "" {
				loggerConfig.Path.Access = filepath.Join(it.Path, loggerConfig.Files.Access)
			}
			if loggerConfig.Path.Error == "" {
				loggerConfig.Path.Error = filepath.Join(it.Path, loggerConfig.Files.Error)
			}
		} else {
			if it.Path != "" { // если путь задан - остальное выставляется по-умолчанию
				if loggerConfig.Files.Access == "" {
					loggerConfig.Files.Access = "access.log"
				}
				if loggerConfig.Files.Error == "" {
					loggerConfig.Files.Error = "error.log"
				}
				if loggerConfig.Path.Access == "" {
					loggerConfig.Path.Access = filepath.Join(it.Path, loggerConfig.Files.Access)
				}
				if loggerConfig.Path.Error == "" {
					loggerConfig.Path.Error = filepath.Join(it.Path, loggerConfig.Files.Error)
				}
			}
		}

		if loggerConfig.Path.Access == "" && it.Path == "" { // пути пусты - пусты и названия файлов
			loggerConfig.Files.Access = ""
		}

		if loggerConfig.Path.Error == "" && it.Path == "" { // пути пусты - пусты и названия файлов
			loggerConfig.Files.Error = ""
		}
	}
}

func (it *ConfigLog) ToLoggersConfigMapFromXMLTemplate(xmlTemplate string) (map[string]string, error) {
	configMap := make(map[string]string)

	t := template.New("log-template")
	if t, e := t.Parse(xmlTemplate); e != nil {
		return nil, fmt.Errorf(`parsing template failed: %s`, e.Error())
	} else {
		var out bytes.Buffer
		for loggerName, loggerConfig := range it.Loggers {
			out.Reset()
			if e := t.Execute(&out, loggerConfig); e != nil {
				return nil, fmt.Errorf(`executing template failed: %s`, e.Error())
			} else {
				configMap[loggerName] = out.String()
			}
		}
	}

	return configMap, nil
}

func (it *ConfigLog) ToLoggersMapFromConfigMap(configMap map[string]string) (*LoggersMap, error) {
	loggersMap := NewLoggersMap()

	for name, config := range configMap {
		if logger, e := it.newLogger(name, config); e != nil {
			return nil, fmt.Errorf("ConfigLog.newLogger('%s', ...) failed: %s", name, e.Error())
		} else {
			loggersMap.Store(name, logger)
		}
	}

	return loggersMap, nil
}

func (it *ConfigLog) ToLoggersMap(xmlTemplatePtr *string) (*LoggersMap, error) {
	xmlTemplate := DefaultXMLTemplate
	if xmlTemplatePtr != nil {
		xmlTemplate = *xmlTemplatePtr
	}
	configMap, e := it.ToLoggersConfigMapFromXMLTemplate(xmlTemplate)
	if e != nil {
		return nil, e
	}

	loggersMap, e := it.ToLoggersMapFromConfigMap(configMap)
	if e != nil {
		return nil, e
	}

	return loggersMap, e
}

func (it *ConfigLog) newLogger(name, config string) (seelog.LoggerInterface, error) {
	logger, e := seelog.LoggerFromConfigAsString(config)
	if e != nil {
		return nil, e
	}

	if name == "app" {
		seelog.ReplaceLogger(logger)
	}

	return logger, nil
}

func (it *ConfigLog) Has(name string) bool {
	_, ok := it.Loggers[name]
	return ok
}

func (it *ConfigLog) Push(name, prefix, min, max string) {
	logger := &configLogger{Prefix: prefix}
	logger.Level.Min = min
	logger.Level.Max = max
	it.Loggers[name] = logger
}

func (it *ConfigLog) Init() {
	if it.Loggers == nil {
		it.Loggers = make(map[string]*configLogger)
	}
}

func NewConfigLog() *ConfigLog {
	it := new(ConfigLog)
	it.Init()
	return it
}
