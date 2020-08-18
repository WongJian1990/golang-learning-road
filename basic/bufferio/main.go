package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("D:/demos/basic/bufferio/readme.txt", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	// defer file.Close()
	r := bufio.NewReader(file)

	//buffered
	count := r.Buffered()
	fmt.Println("count1: ", count)

	//peek
	b, err := r.Peek(10)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("b: ", string(b))

	// err = r.UnreadByte()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	v, err := r.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("v2: %c\n", v)

	err = r.UnreadByte()
	if err != nil {
		log.Fatalln(err)
	}

	v, err = r.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("v2: %c\n", v)

	//read
	var p [10]byte
	len, err := r.Read(p[:])
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("len: ", len, " : ", string(p[:]))
	//buffered
	count = r.Buffered()
	fmt.Println("count2: ", count)

	r.ReadLine()
	line, err := r.ReadSlice('\n')
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("line: ", string(line))
	file.Close()

	file, err = os.OpenFile("readme.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	_, err = w.Write([]byte(string("Test Write API 1\n")))
	if err != nil {
		log.Fatalln(err)
	}
	_, err = w.Write([]byte(string("Test Write API 2\n")))
	if err != nil {
		log.Fatalln(err)
	}
	w.Flush()
	w.Reset(file)
	_, err = w.Write([]byte(string("Test Write API 5\n")))
	if err != nil {
		log.Fatalln(err)
	}
	w.Flush()
}
