package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
	_, err = w.WriteString("Test Write API 6 \n")
	if err != nil {
		log.Fatalln(err)
	}
	w.Flush()

	file.Close()
	//Scanner
	file, err = os.OpenFile("readme.txt", os.O_RDONLY, 0666)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading standard input: %v", err)
	}

	//custom
	// const input = "1234 5678 123456789876543210"
	// scanner = bufio.NewScanner(strings.NewReader(input))
	// split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// 	advance, token, err = bufio.ScanWords(data, atEOF)
	// 	if err == nil && token != nil {
	// 		_, err = strconv.ParseInt(string(token), 10, 32)
	// 	}
	// 	return
	// }
	// scanner.Split(split)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }
	// if err := scanner.Err(); err != nil {
	// 	log.Fatalf("Invalid input: %v", err)
	// }

	//scanner word

	const input2 = "2020/08/19 16:54:36 Invalid input: strconv.ParseInt: parsing \"123456789876543210\": value out of range"
	scanner = bufio.NewScanner(strings.NewReader(input2))
	scanner.Split(bufio.ScanWords)
	count = 0
	for scanner.Scan() {
		count++
		fmt.Println("Text: ", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Invalid input:%v", err)
	}
	fmt.Println("count: ", count)
}
