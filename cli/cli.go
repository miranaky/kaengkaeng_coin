package cli

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/miranaky/kaengkaengcoin/explorer"
	"github.com/miranaky/kaengkaengcoin/rest"
)

func usage() {
	fmt.Printf("Welcome to KaengKaeng Coin!\n\n")
	fmt.Printf("Please use the folllowing flags:\n\n")
	fmt.Printf("	-mode :	Choose between 'html','rest' or 'both'\n")
	fmt.Printf("	-port :	Set port of the server\n\n")
	runtime.Goexit()
}

var port int
var mode string
var help bool

func initFlags() {
	flag.IntVar(&port, "port", 4000, "Set PORT of the server")
	flag.StringVar(&mode, "mode", "rest", "Choose between 'html','rest' or 'both'")
	flag.BoolVar(&help, "help", false, "Show the kanegkaeng coin usage")
}

func Start() {

	initFlags()
	flag.Parse()

	if help {
		usage()
	}

	switch mode {
	case "rest":
		rest.Start(port)
	case "html":
		explorer.Start(port)
	case "both":
		go explorer.Start(port + 1)
		rest.Start(port)
	default:
		usage()
	}
}
