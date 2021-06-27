package main

import (
	"crypto/sha256"
	"fmt"
)

//chain of block

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockChain struct {
	blocks []block
}

func (b *blockChain) getLastHash() string {
	if len(b.blocks) > 0 {
		return b.blocks[len(b.blocks)-1].hash
	}
	return ""
}

func (b *blockChain) addBlock(data string) {
	newBlock := block{data, "", b.getLastHash()}
	hash := sha256.Sum256([]byte(newBlock.data + newBlock.prevHash))
	newBlock.hash = fmt.Sprintf("%x", hash)
	b.blocks = append(b.blocks, newBlock)

}

func (b *blockChain) listBlock() {
	for _, block := range b.blocks {
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Printf("Prev Hash: %s\n", block.prevHash)
	}
}

func main() {
	chain := blockChain{}
	chain.addBlock("Genesis Block")
	chain.addBlock("Second Block")
	chain.addBlock("Third Block")
	chain.listBlock()
}
