package logger

import "github.com/cihub/seelog"

type Getter interface {
	GetLogger() seelog.LoggerInterface
}
