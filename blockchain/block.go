package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

//Block 区块
type Block struct {
	Timestamp int64
	//Data          []byte
	Transactions  []*Transaction
	PrevBlockHash []byte
	BlockHash     []byte
	Nonce         int
}

//NewBlock 创建区块
func NewBlock(transactions []*Transaction /*data string*/, PrevBlockHash []byte) *Block {
	var block *Block
	block = &Block{time.Now().UnixNano(), transactions /*[]byte(data)*/, PrevBlockHash, []byte{}, 0}
	pow := NewProofWork(block)
	nonce, hash := pow.Run()
	block.BlockHash = hash[:]
	block.Nonce = nonce
	return block
}

// //NewGenBlock 创建创世块
// func NewGenBlock() *Block {
// 	return NewBlock("Genesis Block", []byte{})
// }

//Serialize 序列化Block
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return res.Bytes()
}

//Deserialize 反序列化
func (b *Block) Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//HashTransactions 计算区块的所有交易的哈希
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range b.Transactions {
		// txHashes = append(txHashes, tx.ID)
		txHashes = append(txHashes, tx.Hash())
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}
