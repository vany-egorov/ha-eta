package node

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logger() interface{} { return zap.L() }

func initLogger(cfg *config) error {
	logCfg := zap.NewProductionConfig()
	logCfg.Encoding = "console"
	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logCfg.EncoderConfig.CallerKey = ""

	v, err := logCfg.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(v)
	return nil
}
