package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

var subsidy = 50

//TxInput 交易输入
type TxInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

//TxOutput 交易输出
type TxOutput struct {
	Value        int
	ScriptPubKey string
}

//Transaction 交易单元
type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

//NewCoinBaseTx 创建一个coinbase交易
func NewCoinBaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxInput{Txid: []byte{}, Vout: -1, ScriptSig: data}
	txout := TxOutput{subsidy, to}
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()
	return &tx
}

//IsCoinBase 判断是否是coinbase交易
func (tx Transaction) IsCoinBase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && (tx.Vin[0].Vout == -1)
}

//SetID hash交易ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

//CanUnlockOutputWith 签名解锁
func (in *TxInput) CanUnlockOutputWith(unlocking string) bool {
	return in.ScriptSig == unlocking
}

//CanBeUnlockedWith 锁定
func (out *TxOutput) CanBeUnlockedWith(unlocking string) bool {
	return out.ScriptPubKey == unlocking
}

//NewUTXOTransaction 创建一笔新的交易
func NewUTXOTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	//找到足够的未花费输出
	acc, validOuptputs := chain.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("Error:Not enought founds")
	}
	for txid, outs := range validOuptputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TxOutput{amount, to})
	//如果UTXO总数超过所需，则找零
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}
