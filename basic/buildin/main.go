package main

import (
	"fmt"
)

func main() {
	//new 分配内存,返回指针类型,初始值0
	i := new(int)
	*i = 10
	fmt.Printf("%v: %v\n", i, *i)

	//make 构造切片/map/channel
	j := make([]int, 10, 50)
	fmt.Printf("cap: %v\n", cap(j))
	fmt.Printf("addr:%p\n", &j)
	fmt.Println()

	copy(j, []int{10, 20, 30, 40, 50})
	fmt.Printf("cap: %v\n", cap(j))
	fmt.Println()

	j = append(j, 10, 20, 30)

	fmt.Printf("cap: %v\n", cap(j))
	fmt.Println()

	j = []int{10, 20, 30, 40, 50}
	fmt.Printf("%p: %v\n", &j, j)

	c := cap(j)
	fmt.Printf("cap: %v\n", c)

	l := len(j)
	fmt.Printf("len: %v\n", l)

	fmt.Println()

	j = append(j, 10, 20, 30)
	fmt.Printf("%p: %v\n", &j, j)

	fmt.Printf("cap: %v", cap(j))

	m := make(map[string]int)
	m["map"] = 1
	m["test"] = 2
	m["method"] = 3
	fmt.Printf("map: %v\n", m)

	delete(m, "test")
	fmt.Printf("map: %v\n", m)

}
