package servers

import (
	"fmt"
	"strings"
)

type Kind uint8

const (
	KindUnknown Kind = iota
	KindINET
	KindUNIX
)

var kindText = map[Kind]string{
	KindUnknown: "unknown",
	KindINET:    "inet / AF_INET",
	KindUNIX:    "unix / AF_UNIX",
}

func (it Kind) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", it.String())), nil
}
func (it Kind) MarshalYAML() (interface{}, error) { return it.String(), nil }

func (it Kind) String() string { return kindText[it] }

func (it *Kind) UnmarshalJSON(data []byte) error {
	v := strings.Replace(string(data), "\"", "", -1)
	v = strings.TrimSpace(v)
	*it = NewKindFromString(v)
	return nil
}

func (it *Kind) UnmarshalYAML(unmarshal func(interface{}) error) error {
	v := ""
	unmarshal(&v)
	*it = NewKindFromString(v)
	return nil
}

func NewKindFromString(v string) Kind {
	v = strings.TrimSpace(v)
	v = strings.Replace(v, "/", "", -1)
	v = strings.Replace(v, "\"", "", -1)
	v = strings.Replace(v, "'", "", -1)
	v = strings.Replace(v, "_", "-", -1)
	v = strings.Replace(v, " ", "-", -1)
	v = strings.ToLower(v)

	switch v {
	case "inet":
		return KindINET
	case "unix":
		return KindUNIX
	}

	return KindUnknown
}
