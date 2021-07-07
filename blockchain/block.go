package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/miranaky/kaengkaengcoin/db"
	"github.com/miranaky/kaengkaengcoin/utils"
)

type Block struct {
	Hash        string `json:"hash"`
	PrevHash    string `json:"prevHash,omitempty"`
	Height      int    `json:"height"`
	Difficulty  int    `json:"difficulty"`
	Nonce       int    `json:"nonce"`
	TimeStamp   int    `json:"timeStamp"`
	Transaction []*Tx  `json:"transaction"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

var ErrNotFound = errors.New("Block Not Found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.TimeStamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		// fmt.Printf("\n\n\nTarget:%s\nHash:%s\nNonce:%d\n\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(prevHash string, height int) *Block {
	block := &Block{
		Hash:        "",
		PrevHash:    prevHash,
		Height:      height,
		Difficulty:  BlockChain().difficulty(),
		Nonce:       0,
		Transaction: []*Tx{makeCoinbaseTx("sungmuk")},
	}
	block.mine()
	block.persist()
	return block
}
