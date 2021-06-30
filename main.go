package main

import (
	"github.com/miranaky/kaengkaengcoin/explorer"
	"github.com/miranaky/kaengkaengcoin/rest"
)

func main() {
	go explorer.Start(5000)
	rest.Start(4000)

}
