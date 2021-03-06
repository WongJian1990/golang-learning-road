package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var targetBits = 16

//ProofOfWork 工作量证明
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

//NewProofWork 构造ProofWork结构
func NewProofWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	proof := &ProofOfWork{b, target}
	return proof
}

func (pow *ProofOfWork) prepare(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.HashTransactions(),
		[]byte(fmt.Sprintf("%x", pow.block.Timestamp)),
		[]byte(fmt.Sprintf("%x", targetBits)),
		[]byte(fmt.Sprintf("%x", nonce)),
	}, []byte{})
	return data
}

//Run 工作证明
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	// fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < math.MaxInt64 {
		data := pow.prepare(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

//Validate 验证Pow
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepare(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
