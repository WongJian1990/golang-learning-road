package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const (
	dbFile      = "blockchain.db"
	blockBucket = "blocks"
	//创世区块永不修改的数据
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

//BlockChain 区块链
//tip为存储的最新的一个块的哈希
//db 存储数据库连接
// type BlockChain struct {
// 	blocks []*Block
// }
type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

//CreateBlockchain 创建一个新的区块链数据库连接
//address 用来接收挖出创世块的奖励
func CreateBlockchain(address string) *BlockChain {
	// return &BlockChain{[]*Block{NewBlock("Genesis Block", []byte{})}}

	//检查区块链数据库是否存在
	if dbExists() {
		fmt.Println("CreateBlockchain:Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	trans := NewCoinBaseTx(address, genesisCoinbaseData)
	gen := NewBlock([]*Transaction{trans}, []byte{})

	//打开DB文件
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	//开启更新事务
	err = db.Update(func(t *bolt.Tx) error {

		b, err := t.CreateBucket([]byte(blockBucket))
		if err != nil {
			return err
		}
		err = b.Put(gen.BlockHash, gen.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), gen.BlockHash)
		if err != nil {
			return err
		}
		tip = gen.BlockHash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &BlockChain{tip, db}
	//return &BlockChain{[]*Block{NewBlock("Genesis Block", []byte{})}}
}

//NewBlockChain 创建一个有创世块的新区块链
func NewBlockChain() *BlockChain {
	if dbExists() == false {
		fmt.Println("No existing block chain found. Create genenis first.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blockBucket))
		if b == nil {
			return errors.New("Get block bucket failed")
		}
		tip = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &BlockChain{tip, db}
}

//MineBlock 挖矿(区块)
func (chain *BlockChain) MineBlock(transactions []*Transaction) *Block {

	if dbExists() == false {
		fmt.Println("MinBlock: block chain not existing. create genesis first.")
		os.Exit(1)
	}
	var lastHash []byte
	for _, tx := range transactions {
		if chain.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}
	//首先获取最新一次的哈希用于生成新的区块的哈希
	//开启一个只读的事务操作
	err := chain.db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blockBucket))
		if b == nil {
			return errors.New("AddBlock: Get block bucket failed")
		}
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(transactions, lastHash)

	err = chain.db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blockBucket))
		if b == nil {
			return errors.New("AddBlock: Get block bucket failed")
		}
		err := b.Put(newBlock.BlockHash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), newBlock.BlockHash)
		if err != nil {
			return err
		}
		chain.tip = newBlock.BlockHash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return newBlock
}

//SignTransaction 对交易输入进行签名
func (chain *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTx, err := chain.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	tx.Sign(privKey, prevTXs)
}

//VerifyTransaction 验证交易有效性
func (chain *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}
	prevTXs := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		prevTx, err := chain.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	return tx.Verify(prevTXs)
}

//BlockChainIterator 区块链迭代器
type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

//Iterator 获取区块链迭代器
func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{
		chain.tip,
		chain.db,
	}
}

//Next 返回区块链的下一个区块
func (it *BlockChainIterator) Next() *Block {
	if len(it.currentHash) == 0 {
		return nil
	}
	var block = &Block{}
	err := it.db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blockBucket))
		if b == nil {
			return errors.New("Next: Get block bucket failed")
		}
		encodeBlock := b.Get(it.currentHash)
		block = block.Deserialize(encodeBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	it.currentHash = block.PrevBlockHash
	return block
}

//FindUTXO 查找交易未花费
func (chain *BlockChain) FindUTXO() map[string]TxOutputs {
	var UTXO = make(map[string]TxOutputs)
	spendTXOs := make(map[string][]int)
	it := chain.Iterator()
	for {
		if block := it.Next(); block != nil {
			for _, tx := range block.Transactions {
				txID := hex.EncodeToString(tx.ID)
			Outputs:
				for outIdx, out := range tx.Vout {
					if spendTXOs[txID] != nil {
						for _, spnedOutIdx := range spendTXOs[txID] {
							if spnedOutIdx == outIdx {
								continue Outputs
							}
						}
					}
					outs := UTXO[txID]
					outs.Outputs = append(outs.Outputs, out)
					UTXO[txID] = outs
				}
				if tx.IsCoinBase() == false {
					for _, in := range tx.Vin {
						inTxID := hex.EncodeToString(in.Txid)
						spendTXOs[inTxID] = append(spendTXOs[inTxID], in.Vout)
					}
				}
			}
		} else {
			break
		}
	}
	return UTXO
}

//FindTransaction 通过ID查找交易
func (chain *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	it := chain.Iterator()
	for {
		if block := it.Next(); block != nil {
			for _, tx := range block.Transactions {
				if bytes.Compare(tx.ID, ID) == 0 {
					return *tx, nil
				}
			}
		} else {
			break
		}
	}
	return Transaction{}, errors.New("Transaction not found")
}
