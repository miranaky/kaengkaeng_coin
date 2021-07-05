package main

import (
	"github.com/miranaky/kaengkaengcoin/cli"
	"github.com/miranaky/kaengkaengcoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
