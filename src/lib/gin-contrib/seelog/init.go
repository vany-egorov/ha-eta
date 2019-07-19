package seelog

import "github.com/cihub/seelog"

type LoggerGetter interface {
	GetLogger() seelog.LoggerInterface
}
