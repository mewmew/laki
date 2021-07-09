package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/mewkiz/pkg/term"
	"github.com/mewmew/laki/vk"
)

var (
	// dbg is a logger with the "laki:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.MagentaBold("laki:")+" ", 0)
	// warn is a logger with the "laki:" prefix which logs warning messages to
	// standard error.
	warn = log.New(os.Stderr, term.RedBold("laki:")+" ", log.Lshortfile)
)

func init() {
	if !debug {
		dbg.SetOutput(ioutil.Discard)
	}
}

// Enable debug output.
const debug = true

func main() {
	if err := vk.Init(); err != nil {
		warn.Fatalf("%+v", err)
	}
}
