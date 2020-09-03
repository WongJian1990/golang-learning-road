package main

import (
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

const (
	dbFile      = "blockchain.db"
	blockBucket = "blocks"
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

//NewBlockChain 创建区块链
func NewBlockChain() *BlockChain {
	// return &BlockChain{[]*Block{NewBlock("Genesis Block", []byte{})}}
	var tip []byte
	//打开DB文件
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	//开启更新事务
	err = db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blockBucket))
		if b == nil {
			//数据库中不存在区块链则创建
			gen := NewBlock("Genesis Block", []byte{})
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
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &BlockChain{tip, db}
	//return &BlockChain{[]*Block{NewBlock("Genesis Block", []byte{})}}
}

//AddBlock 添加区块
func (chain *BlockChain) AddBlock(data string) {
	// prevBlock := chain.blocks[len(chain.blocks)-1]
	// block := NewBlock(data, prevBlock.BlockHash)
	// chain.blocks = append(chain.blocks, block)

	var lastHash []byte
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
	newBlock := NewBlock(data, lastHash)

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
