package environment

import (
	"fmt"
	"strings"
)

type Environment uint8

const (
	EnvUnknown Environment = 0

	EnvDev Environment = iota
	EnvTest
	EnvDemo
	EnvProd
)

var environmentText = map[Environment]string{
	EnvUnknown: "unknown",
	EnvDev:     "development",
	EnvTest:    "test",
	EnvDemo:    "demo",
	EnvProd:    "production",
}

func EnvironmentValidList() (validList []string) {
	for m, name := range environmentText {
		if m.IsUnknown() {
			continue
		}
		validList = append(validList, name)
	}

	return
}

func EnvironmentValidListAsString() string { return strings.Join(EnvironmentValidList(), " | ") }

func NewEnvironment(v string) Environment {
	v = strings.ToLower(v)
	v = strings.Replace(v, " ", "", -1)
	v = strings.Replace(v, "_", "-", -1)

	switch v {
	case "d", "dev", "development":
		return EnvDev
	case "t", "tst", "test":
		return EnvTest
	case "demo":
		return EnvDemo
	case "p", "prod", "production":
		return EnvProd
	}

	return EnvUnknown
}

func (self Environment) Is(v Environment) bool { return self == v }

func (self Environment) IsDev() bool     { return self.Is(EnvDev) }
func (self Environment) IsTest() bool    { return self.Is(EnvTest) }
func (self Environment) IsNotTest() bool { return !self.Is(EnvTest) }
func (self Environment) IsDemo() bool    { return self.Is(EnvDemo) }
func (self Environment) IsProd() bool    { return self.Is(EnvProd) }
func (self Environment) IsUnknown() bool { return self.Is(EnvUnknown) }
func (self Environment) IsValid() bool   { return !self.IsUnknown() }

func (self Environment) String() string                  { return environmentText[self] }
func (self Environment) ValidList() (validList []string) { return EnvironmentValidList() }
func (self Environment) ValidListAsString() string       { return EnvironmentValidListAsString() }
func (self Environment) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", self.String())), nil
}

func (self *Environment) UnmarshalJSON(data []byte) error {
	v := strings.Replace(string(data), "\"", "", -1)
	*self = NewEnvironment(v)
	return nil
}

func (self *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	v := ""
	unmarshal(&v)
	*self = NewEnvironment(v)
	if self.IsUnknown() {
		return fmt.Errorf("got unknown environment '%s', possible environments are number of: %s;", v, EnvironmentValidListAsString())
	}
	return nil
}
