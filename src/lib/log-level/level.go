package logLvl

type Level uint8

const (
	Quiet Level = iota
	Trace
	Debug
	Info
	Warn
	Error
	Critical
)
