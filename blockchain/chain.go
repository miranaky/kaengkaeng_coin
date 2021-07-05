package blockchain

import (
	"fmt"
	"sync"

	"github.com/miranaky/kaengkaengcoin/db"
	"github.com/miranaky/kaengkaengcoin/utils"
)

type blockChain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockChain
var once sync.Once

func (b *blockChain) persist() {
	db.SaveBlockChain(utils.ToBytes(b))
}

func (b *blockChain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()

}

func (b *blockChain) restore(data []byte) {
	utils.FromBytes(b, data)
}

// GetBlockChain first initialized the blockChain struct by singleton pattern.
func BlockChain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{"", 0}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis Block")
			} else {
				b.restore(checkpoint)
			}
		})
	}
	fmt.Println(b.NewestHash)
	return b
}
