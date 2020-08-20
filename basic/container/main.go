package main

import (
	"container/heap"
	"container/list"
	"fmt"
)

type Node struct {
	value    string
	priority int
	index    int
}

type PriorityQueue []*Node

func (p PriorityQueue) Len() int {
	return len(p)
}

func (p PriorityQueue) Less(i, j int) bool {
	return p[i].priority > p[j].priority
}

func (p PriorityQueue) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

func (p *PriorityQueue) Push(x interface{}) {
	n := len(*p)
	node := x.(*Node)
	node.index = n
	*p = append(*p, node)
}

func (p *PriorityQueue) Pop() interface{} {
	old := *p
	n := len(old)
	node := old[n-1]
	node.index = -1
	*p = old[0 : n-1]
	return node
}

func (p *PriorityQueue) update(node *Node, value string, priority int) {
	node.value = value
	node.priority = priority
	heap.Fix(p, node.index)
}

func main() {
	items := map[string]int{
		"banana": 3,
		"apple":  2,
		"pear":   4,
	}

	p := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		p[i] = &Node{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&p)

	item := &Node{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&p, item)
	p.update(item, item.value, 5)
	for p.Len() > 0 {
		item := heap.Pop(&p).(*Node)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
	fmt.Println()

	lst := list.New()
	len := lst.Len()
	fmt.Printf("len: %v\n", len)
	lst.PushBack("1")
	lst.PushBack("2")
	lst.PushBack("3")
	el := lst.Front()
	for el != nil {
		fmt.Printf("%v ", el.Value.(string))
		el = el.Next()

	}
	fmt.Println()
}
