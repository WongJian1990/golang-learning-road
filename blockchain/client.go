package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

//CLI 区块链客户端操作结构
type CLI struct {
	chain *BlockChain
}

const usage = `
Usage:
	add -data BLOCK_DATA add a block to blockchain
	print 	print all blocks in the blockchain
`

func NewClient(chain *BlockChain) *CLI {
	if chain == nil {
		return nil
	}
	return &CLI{chain}
}
func (c *CLI) printUsage() {
	fmt.Println(usage)
}

func (c *CLI) validateArgs() {
	if len(os.Args) < 2 {
		c.printUsage()
		os.Exit(1)
	}
}

func (c *CLI) add(data string) {
	c.chain.AddBlock(data)
}
func (c *CLI) print() {
	it := c.chain.Iterator()
	for {
		if block := it.Next(); block != nil {
			fmt.Printf("Prev Hash: %x\n", block.PrevBlockHash)
			fmt.Printf("Data:%s\n", block.Data)
			fmt.Printf("Current Hash: %x\n", block.BlockHash)
			pow := NewProofWork(block)
			fmt.Printf("PoW:%s\n", strconv.FormatBool(pow.Validate()))
			fmt.Println()
		} else {
			break
		}
	}

}

//Run 运行客户端
func (c *CLI) Run() {
	c.validateArgs()
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printBlockCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "block data")
	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print":
		err := printBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		c.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		//c.chain.AddBlock(*addBlockData)
		c.add(*addBlockData)
	}

	if printBlockCmd.Parsed() {
		c.print()
	}

}
