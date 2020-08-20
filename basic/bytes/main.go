package main

import (
	"bytes"
	"fmt"
	"unicode"
)

func main() {
	a := []byte("Hello, CPlusplus")
	b := []byte("Hello, World")
	c := []byte("Hello, Golang")
	d := a
	res := bytes.Compare(a, b)
	fmt.Printf("res: %v\n", res)
	res = bytes.Compare(b, c)
	fmt.Printf("res: %v\n", res)
	res = bytes.Compare(a, d)
	fmt.Printf("res: %v\n", res)

	fmt.Printf("equal: %v\n", bytes.Equal(a, b))
	fmt.Printf("equal: %v\n", bytes.Equal(a, d))

	rune := bytes.Runes([]byte("这是Unicode 么"))
	fmt.Printf("rune: %v\n", rune)

	fmt.Printf("hasPrefix: %v\n", bytes.HasPrefix([]byte("/usr/share"), []byte("/")))

	fmt.Printf("hasSuffix: %v\n", bytes.HasSuffix([]byte("index.html"), []byte(".html")))

	fmt.Printf("contains: %v\n", bytes.Contains([]byte("Hello, Golang"), []byte("Golang")))

	fmt.Printf("count: %v\n", bytes.Count([]byte("Hello, Golang, Golang, Golang"), []byte("Golang")))

	fmt.Printf("index: %v\n", bytes.Index([]byte("Hello, Golang, Golang, Golang"), []byte("Golang")))

	fmt.Printf("indexByte: %v\n", bytes.IndexByte([]byte("Hello, Golang"), 'G'))

	fmt.Printf("indexAny: %v\n", bytes.IndexAny([]byte("Hello, Golang"), "oxz"))
	f := func(r int32) bool {
		if r == int32('n') {
			return true
		}
		return false
	}
	fmt.Printf("indexFunc: %v\n", bytes.IndexFunc([]byte("Hello, Golang"), f))

	fmt.Printf("lastIndex: %v\n", bytes.LastIndex([]byte("Hello, Golang, Golang, Golang"), []byte("Golang")))

	fmt.Printf("toLower: %v\n", string(bytes.ToLower([]byte("Hello, Golang"))))

	fmt.Printf("toLower: %v\n", string(bytes.ToLowerSpecial(unicode.TurkishCase, []byte("1234 Hello, Golang"))))

	fmt.Printf("repeat: %v\n", string(bytes.Repeat([]byte("Hello, Golang "), 3)))

	fmt.Printf("replace: %v\n", string(bytes.Replace([]byte("Hello, Golang, Golang, Golang"), []byte("Golang"), []byte("CPlusplus"), 1)))
	fmt.Printf("replace: %v\n", string(bytes.Replace([]byte("Hello, Golang, Golang, Golang"), []byte("Golang"), []byte("CPlusplus"), 2)))

	fmt.Printf("trim: %v\n", string(bytes.Trim([]byte(",Hello, Golang, Golang, Golang,"), ",")))
	fmt.Printf("trimSpace: %v\n", string(bytes.TrimSpace([]byte(" Hello Golang Golang Golang "))))

	fmt.Printf("trimFunc: %v\n", string(bytes.TrimFunc([]byte("z Hello Golang k"), func(r int32) bool {
		if r > int32('j') {
			return true
		}
		return false
	})))

	fmt.Printf("TrimLeft: %v\n", string(bytes.TrimLeft([]byte("zk Hello Golang"), "zk")))

	fmt.Printf("TrimLeftFunc: %v\n", string(bytes.TrimLeftFunc([]byte("efghi Hello Golang"), func(r int32) bool {
		if r >= int32('e') {
			return true
		}
		return false
	})))

	slice := bytes.Fields([]byte("Hello World By Golang"))
	for i, v := range slice {
		fmt.Printf("fields[%v]: %v ", i, string(v))
	}
	fmt.Println()

	slice = bytes.FieldsFunc([]byte("Hello, World; By, DK"), func(r int32) bool {
		if r == int32(',') || r == int32(';') || r == int32(' ') {
			return true
		}
		return false
	})

	for i, v := range slice {
		fmt.Printf("fields[%v]: %v ", i, string(v))
	}
	fmt.Println()

	slice = bytes.Split([]byte("Good Boy !!!"), []byte(" "))

	for i, v := range slice {
		fmt.Printf("fields[%v]: %v ", i, string(v))
	}
	fmt.Println()

	slice = bytes.SplitAfter([]byte("Good Boy !!!"), []byte(" "))

	for i, v := range slice {
		fmt.Printf("fields[%v]: %v ", i, string(v))
	}
	fmt.Println()

	fmt.Printf("join: %v\n", string(bytes.Join([][]byte{[]byte("News "), []byte("BBC")}, []byte(", "))))

	r := bytes.NewReader([]byte("Hello, Golang"))
	fmt.Printf("len: %v\n", r.Len())
	//...

	//buffer
	buf := bytes.NewBuffer([]byte("Hello, Golang"))
	fmt.Printf("len: %v\n", buf.Len())
	fmt.Printf("next: %v\n", string(buf.Next(5)))
	fmt.Printf("bytes: %v\n", string(buf.Bytes()))
	fmt.Printf("string: %v\n", buf.String())

	buf.Truncate(4)
	fmt.Printf("truncate: %v\n", buf.String())
}
