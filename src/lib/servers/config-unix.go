package servers

import "os"

type ConfigUNIX struct {
	Addr           string      `yaml:"addr"`
	SocketFileMode os.FileMode `yaml:"socket-file-mode"`
}

func (it *ConfigUNIX) GetAddr() string { return it.Addr }

func (it *ConfigUNIX) ToLog(log ConfigLogger) {
	f := log.Debugf

	f("      addr: %s", it.Addr)
	f("      socket-file-mode: %03o | %s", it.SocketFileMode, it.SocketFileMode)
}

func (it *ConfigUNIX) Defaultize() error {
	if it.SocketFileMode == 0 {
		it.SocketFileMode = DefaultUNIXSocketFileMode | os.ModeSocket
	}

	return nil
}
