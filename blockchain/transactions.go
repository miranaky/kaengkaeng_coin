package blockchain

import (
	"errors"
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

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

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

func makeTx(from, to string, amount int) (*Tx, error) {
	if BlockChain().BalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	oldTxOuts := BlockChain().TxOutsByAddress(from)
	for _, txOut := range oldTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{txOut.Onwer, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txOut.Amount
	}
	change := total - amount
	if change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id:        "",
		TimeStamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("sungmuk", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("sungmuk")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
