package blockchain

import (
	"time"

	"github.com/miranaky/kaengkaengcoin/utils"
)

const (
	minerReward int = 50
)

type Tx struct {
	Id        string   `json:"id"`
	TimeStamp int      `json:"timeStamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	Onwer  string `json:"onwer"`
	Amount int    `json:"amount"`
}

type TxOut struct {
	Onwer  string `json:"onwer"`
	Amount int    `json:"amount"`
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"COINBASE", minerReward},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		TimeStamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}
