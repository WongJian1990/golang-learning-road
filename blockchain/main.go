package main

func main() {
	chain := NewBlockChain()
	defer chain.db.Close()
	cli := NewClient(chain)
	if cli != nil {
		cli.Run()
	}
}
