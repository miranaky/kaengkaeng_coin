package main

import (
	"github.com/miranaky/kaengkaengcoin/blockchain"
	"github.com/miranaky/kaengkaengcoin/cli"
)

func main() {
	blockchain.BlockChain()
	cli.Start()
}
