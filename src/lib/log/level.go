package log

type Level uint8

const (
	/// A special log level used to turn off logging.
	Quiet Level = iota

	// For pervasive information on states of all elementary constructs.
	// Use 'Trace' for in-depth debugging to find problem parts of a function,
	// to check values of temporary variables, etc.
	Trace

	// For detailed system behavior reports and diagnostic messages
	// to help to locate problems during development.
	Debug

	// For general information on the application's work.
	// Use 'Info' level in your code so that you could leave it
	// 'enabled' even in production. So it is a 'production log level'.
	Info

	// For indicating small errors, strange situations,
	// failures that are automatically handled in a safe manner.
	Warn

	// For severe failures that affects application's workflow,
	// not fatal, however (without forcing app shutdown).
	Error

	// For producing final messages before application’s death.
	Critical
)
