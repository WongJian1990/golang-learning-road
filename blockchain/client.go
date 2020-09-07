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
	balance 			-address ADDRESS - Get balance of ADDRESS
	createblockchain 	-adress	 ADDRESS - Create a blockchain and send genesis block reward to ADDRESS
	createwallet		- Generates a new key-pair and saves it into the wallet file
	send				-from 	 FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO address
	reindexutxo 		- Rebuilds the UTXO set
	listaddress			- Lists all addresses from wallet file	
	print 				- Print all blocks in the blockchain
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
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	chain := CreateBlockchain(address)
	defer chain.db.Close()
	UTXOSet := UTXOSet{chain}
	UTXOSet.Reindex()
	fmt.Println("Done!")
}

func (c *CLI) createWallets() {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	fmt.Printf("Your new address:%s\n", address)
}

func (c *CLI) print() {
	chain := NewBlockChain()
	defer chain.db.Close()
	it := chain.Iterator()
	for {
		if block := it.Next(); block != nil {
			fmt.Printf("=====Block: %x======\n", block.BlockHash)
			fmt.Printf("Prev Hash: %x\n", block.PrevBlockHash)
			pow := NewProofWork(block)
			fmt.Printf("PoW:%s\n", strconv.FormatBool(pow.Validate()))
			fmt.Println()
		} else {
			break
		}
	}
}

func (c *CLI) listAddresses() {
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()
	for _, addr := range addresses {
		fmt.Println(addr)
	}
}

func (c *CLI) send(from, to string, amount int) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}
	chain := NewBlockChain()
	UTXOSet := UTXOSet{chain}
	defer chain.db.Close()
	tx := NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := NewCoinBaseTx(from, "")
	txs := []*Transaction{cbTx, tx}
	newBlock := chain.MineBlock(txs)
	UTXOSet.Update(newBlock)

	fmt.Println("Success!")
}

func (c *CLI) reindexUTXO() {
	chain := NewBlockChain()
	UTXOSet := UTXOSet{chain}
	UTXOSet.Reindex()
	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}

//balance 账号余额
func (c *CLI) balance(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	chain := NewBlockChain()
	UTXOSet := UTXOSet{chain}
	defer chain.db.Close()
	amount := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)
	for _, out := range UTXOs.Outputs {
		amount += out.Value
	}
	fmt.Printf("Blance of '%s':%d\n", address, amount)
}

//Run 运行客户端
func (c *CLI) Run() {
	c.validateArgs()
	blanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("cratewallet", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listaddress", flag.ExitOnError)
	printBlockCmd := flag.NewFlagSet("print", flag.ExitOnError)

	balanceAddress := blanceCmd.String("address", "", "the address of balance")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "the address to send genesis block reword to")
	sendFrom := sendCmd.String("from", "", "source wallet address")
	sendTo := sendCmd.String("to", "", "destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "amount to send")

	switch os.Args[1] {
	case "blance":
		err := blanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddress":
		err := listAddressCmd.Parse(os.Args[2:])
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
	if createBlockchainCmd.Parsed() {
		c.createBlockchain(*createBlockchainAddress)
	}
	if createWalletCmd.Parsed() {
		c.createWallets()
	}

	if listAddressCmd.Parsed() {
		c.listAddresses()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		c.send(*sendFrom, *sendTo, *sendAmount)
	}

	if reindexUTXOCmd.Parsed() {
		c.reindexUTXO()
	}
	if printBlockCmd.Parsed() {
		c.print()
	}

}
