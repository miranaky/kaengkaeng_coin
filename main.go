package main

import (
	"fmt"

	"github.com/miranaky/kaengkaengcoin/blockchain"
)

func main() {
	chain := blockchain.GetBlockChain()
	chain.AddBlock("Second block")
	chain.AddBlock("Thrid block")
	chain.AddBlock("Fourth block")
	for _, block := range chain.AllBlocks() {
		fmt.Printf("Data: %s \n", block.Data)
		fmt.Printf("Hash: %s \n", block.Hash)
		fmt.Printf("Prev Hash: %s \n", block.PrevHash)

	}
}
