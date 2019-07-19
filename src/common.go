package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/vany-egorov/ha-eta/lib"
)

var (
	buildDate string
	version   string

	doneAt *lib.DoneAt = lib.NewDefaultDoneAt()
)

func initialize() error {
	rand.Seed(time.Now().UnixNano())

	if it, e := lib.NewDoneAt(buildDate); e != nil {
		return e
	} else {
		doneAt = it
	}

	return nil
}

func MustInitialize() {
	if e := initialize(); e != nil {
		fmt.Fprintf(os.Stderr, "initialization failed: %s", e.Error())
		os.Exit(1)
	}
}
