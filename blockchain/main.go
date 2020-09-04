package main

func main() {

	cli := NewClient()
	if cli != nil {
		cli.Run()
	}
}
