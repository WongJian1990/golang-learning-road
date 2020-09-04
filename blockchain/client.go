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
	// chain *BlockChain
}

const usage = `
Usage:
	balance -address ADDRESS - Get balance of ADDRESS
	create 	-adress	 ADDRESS - Create a blockchain and send genesis block reward to ADDRESS
	send	-from 	 FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO address
	print 	print all blocks in the blockchain
`

//NewClient 客户端
func NewClient() *CLI {

	return &CLI{}
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

// func (c *CLI) add(data string) {
// 	c.chain.AddBlock(data)
// }

func (c *CLI) createBlockchain(address string) {
	chain := CreateBlockchain(address)
	defer chain.db.Close()
	fmt.Println("Done!")
}

func (c *CLI) print() {
	chain := NewBlockChain("")
	defer chain.db.Close()
	it := chain.Iterator()
	for {
		if block := it.Next(); block != nil {
			fmt.Printf("Prev Hash: %x\n", block.PrevBlockHash)
			// fmt.Printf("Data:%s\n", block.Data)
			fmt.Printf("Current Hash: %x\n", block.BlockHash)
			pow := NewProofWork(block)
			fmt.Printf("PoW:%s\n", strconv.FormatBool(pow.Validate()))
			fmt.Println()
		} else {
			break
		}
	}
}

func (c *CLI) send(from, to string, amount int) {
	chain := NewBlockChain(from)
	defer chain.db.Close()
	tx := NewUTXOTransaction(from, to, amount, chain)
	chain.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}

//balance 账号余额
func (c *CLI) balance(address string) {
	chain := NewBlockChain(address)
	defer chain.db.Close()
	amount := 0
	UTXOs := chain.FindUTXO(address)
	for _, out := range UTXOs {
		amount += out.Value
	}
	fmt.Printf("Blance of '%s':%d\n", address, amount)
}

//Run 运行客户端
func (c *CLI) Run() {
	c.validateArgs()
	blanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printBlockCmd := flag.NewFlagSet("print", flag.ExitOnError)

	balanceAddress := blanceCmd.String("address", "", "the address of balance")
	createAddress := createCmd.String("address", "", "the address to send genesis block reword to")
	sendFrom := sendCmd.String("from", "", "source wallet address")
	sendTo := sendCmd.String("to", "", "destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "amount to send")

	switch os.Args[1] {
	case "blance":
		err := blanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "create":
		err := createCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
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
	if blanceCmd.Parsed() {
		c.balance(*balanceAddress)
	}
	if createCmd.Parsed() {
		c.createBlockchain(*createAddress)
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		c.send(*sendFrom, *sendTo, *sendAmount)
	}
	if printBlockCmd.Parsed() {
		c.print()
	}

}
