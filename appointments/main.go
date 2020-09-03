package main

import (
	"fmt"
	"sync"
	"time"
)

// import (
// 	"fmt"
// 	"runtime"
// 	"sync"
// )

// type student struct {
// 	Name string
// 	Age  int
// }

// func pase_student() {
// 	m := make(map[string]*student)
// 	stus := []student{
// 		{Name: "zhou", Age: 24},
// 		{Name: "li", Age: 23},
// 		{Name: "wang", Age: 22},
// 	}
// 	for _, stu := range stus {
// 		m[stu.Name] = &stu
// 	}
// 	for _, stu := range m {
// 		println("stu: ", stu.Name, stu.Age)
// 	}
// }

// func main() {
// 	//pase_student()
// 	//defer_call()

// 	runtime.GOMAXPROCS(1)
// 	wg := sync.WaitGroup{}
// 	wg.Add(20)
// 	for i := 0; i < 10; i++ {
// 		go func() {
// 			fmt.Println("i: ", i)
// 			wg.Done()
// 		}()
// 	}
// 	// runtime.Gosched()

// 	for i := 0; i < 10; i++ {
// 		go func(i int) {
// 			fmt.Println("i: ", i)
// 			wg.Done()
// 		}(i)
// 	}

// 	wg.Wait()
// }

// func defer_call() {
// 	defer func() { fmt.Println("打印前") }()
// 	defer func() { fmt.Println("打印中") }()
// 	defer func() { fmt.Println("打印后") }()

// 	//panic("触发异常")
// }

// type People struct{}

// func (p *People) ShowA() {
// 	fmt.Println("showA")
// 	p.ShowB()
// }
// func (p *People) ShowB() {
// 	fmt.Println("showB")
// }

// type Teacher struct {
// 	People
// }

// func (t *Teacher) ShowB() {
// 	fmt.Println("teacher showB")
// }

// func main() {
// 	t := Teacher{}
// 	t.ShowA()
// }

// func main() {
// 	runtime.GOMAXPROCS(1)
// 	int_chan := make(chan int, 1)
// 	string_chan := make(chan string, 1)
// 	int_chan <- 1
// 	string_chan <- "hello"
// 	select {
// 	case value := <-int_chan:
// 		fmt.Println(value)
// 	case value := <-string_chan:
// 		panic(value)
// 	}
// }

// func calc(index string, a, b int) int {
// 	ret := a + b
// 	fmt.Println(index, a, b, ret)
// 	return ret
// }

// func main() {
// 	a := 1
// 	b := 2
// 	defer calc("1", a, calc("10", a, b))
// 	a = 0
// 	defer calc("2", a, calc("20", a, b))
// 	b = 1
// }

// type UserAges struct {
// 	ages map[string]int
// 	sync.Mutex
// }

// func (ua *UserAges) Add(name string, age int) {
// 	ua.Lock()
// 	defer ua.Unlock()
// 	ua.ages[name] = age
// }

// func (ua *UserAges) Get(name string) int {
// 	if age, ok := ua.ages[name]; ok {
// 		return age
// 	}
// 	return -1
// }

// func main() {
// 	s := make([]int, 5)
// 	s = append(s, 1, 2, 3)
// 	fmt.Println(s)
// 	user := &UserAges{ages: make(map[string]int)}
// 	go user.Add("lihua", 28)
// 	go user.Get("lihua")
// 	// runtime.Gosched()
// }

// type People interface {
// 	Speak(string) string
// }

// type Student struct{}

// func (stu *Student) Speak(think string) (talk string) {
// 	if think == "bitch" {
// 		talk = "You are a good boy"
// 	} else {
// 		talk = "hi"
// 	}
// 	return
// }

// func main() {
// 	var peo People = &Student{}
// 	think := "bitch"
// 	fmt.Println(peo.Speak(think))
// }

// type People interface {
// 	Show()
// }

// type Student struct{}

// func (stu *Student) Show() {

// }

// func live() People {
// 	var stu *Student
// 	return stu
// }

// func main() {
// 	if live() == nil {
// 		fmt.Println("AAAAAAA")
// 	} else {
// 		fmt.Println("BBBBBBB")
// 	}
// }

// func main() {
// 	i := GetValue()

// 	switch i.(type) {
// 	case int:
// 		println("int")
// 	case string:
// 		println("string")
// 	case interface{}:
// 		println("interface")
// 	default:
// 		println("unknown")
// 	}

// }

// func GetValue() int {
// 	return 1
// }

// func funcMui(x,y int)(sum int,error){
//     return x+y,nil
// }

// func main() {

// 	println(DeferFunc1(1))
// 	println(DeferFunc2(1))
// 	println(DeferFunc3(1))
// }

// func DeferFunc1(i int) (t int) {
// 	t = i
// 	defer func() {
// 		t += 3
// 	}()
// 	return t
// }

// func DeferFunc2(i int) int {
// 	t := i
// 	defer func() {
// 		t += 3
// 	}()
// 	return t
// }

// func DeferFunc3(i int) (t int) {
// 	defer func() {
// 		t += i
// 	}()
// 	return 2
// }

// func main() {
// 	list := new([]int)
// 	list = append(list, 1)
// 	fmt.Println(list)
// }

// func main() {
// 	s1 := []int{1, 2, 3}
// 	s2 := []int{4, 5}
// 	s1 = append(s1, s2)
// 	fmt.Println(s1)
// }

// func main() {

// 	sn1 := struct {
// 		age  int
// 		name string
// 	}{age: 11, name: "qq"}
// 	sn2 := struct {
// 		age  int
// 		name string
// 	}{age: 11, name: "qq"}

// 	if sn1 == sn2 {
// 		fmt.Println("sn1 == sn2")
// 	}

// 	// sm1 := struct {
// 	// 	age int
// 	// 	m   map[string]string
// 	// }{age: 11, m: map[string]string{"a": "1"}}
// 	// sm2 := struct {
// 	// 	age int
// 	// 	m   map[string]string
// 	// }{age: 11, m: map[string]string{"a": "1"}}

// 	// if sm1 == sm2 {
// 	// 	fmt.Println("sm1 == sm2")
// 	// }
// }

// func Foo(x interface{}) {
// 	if x == nil {
// 		fmt.Println("empty interface")
// 		return
// 	}
// 	fmt.Println("non-empty interface")
// }
// func main() {
// 	var x *int = nil
// 	Foo(x)
// }

// func GetValue(m map[int]string, id int) (string, bool) {
// 	if _, exist := m[id]; exist {
// 		return "存在数据", true
// 	}
// 	return nil, false
// }
// func main() {
// 	intmap := map[int]string{
// 		1: "a",
// 		2: "bb",
// 		3: "ccc",
// 	}

// 	v, err := GetValue(intmap, 3)
// 	fmt.Println(v, err)
// }

// const (
// 	x = 10
// 	y = iota * iota
// 	z
// 	_
// 	k
// 	p = iota
// )

// func main() {
// 	fmt.Println(x, y, z, k, p)
// }

// type Param map[string]interface{}

// type Show struct {
// 	Param
// }

// func main() {
// 	s := new(Show)
// 	s.Param["RMB"] = 10000
// }

// type People struct {
// 	name string `json:"name"`
// }

// func main() {
// 	js := `{
// 		"name":"11"
// 	}`
// 	var p People
// 	err := json.Unmarshal([]byte(js), &p)
// 	if err != nil {
// 		fmt.Println("err: ", err)
// 		return
// 	}
// 	fmt.Println("people: ", p)
// }

// type People struct {
// 	Name string
// }

// func (p *People) String() string {
// 	return fmt.Sprintf("print: %v", p)
// }

// func main() {
// 	p := &People{}
// 	p.String()
// }

// func main() {
// 	ch := make(chan int, 1000)
// 	go func() {
// 		for i := 0; i < 10; i++ {
// 			ch <- i
// 		}
// 	}()
// 	go func() {
// 		for {
// 			a, ok := <-ch
// 			if !ok {
// 				fmt.Println("close")
// 				return
// 			}
// 			fmt.Println("a: ", a)
// 		}
// 	}()
// 	// runtime.Gosched()
// 	close(ch)
// 	fmt.Println("ok")
// 	time.Sleep(time.Second * 100)
// }

// type Project struct{}

// func (p *Project) deferError() {
// 	if err := recover(); err != nil {
// 		fmt.Println("recover: ", err)
// 	}
// }

// func (p *Project) exec(msgchan chan interface{}) {
// 	for msg := range msgchan {
// 		m := msg.(int)
// 		fmt.Println("msg: ", m)
// 	}
// }

// func (p *Project) run(msgchan chan interface{}) {
// 	for {
// 		defer p.deferError()
// 		go p.exec(msgchan)
// 		time.Sleep(time.Second * 2)
// 	}
// }

// func (p *Project) Main() {
// 	a := make(chan interface{}, 100)
// 	go p.run(a)
// 	go func() {
// 		for {
// 			a <- "1"
// 			time.Sleep(time.Second)
// 		}
// 	}()
// 	time.Sleep(time.Second * 1000000000)
// }

// func main() {
// 	p := new(Project)
// 	p.Main()
// }

// func Foo(x interface{}) {
// 	if x == nil {
// 		fmt.Println("empty interface")
// 		return
// 	}
// 	fmt.Println("non-empty interface")
// }
// func main() {
// 	var x *int = nil
// 	Foo(x)
// }

// type Student struct {
// 	Name string
// }

// func main() {
// 	m := map[string]Student{"people": {"zhoujielun"}}
// 	m["people"].Name = "wuyanzu"
// }

// type query func(string) string

// func exec(name string, vs ...query) string {
// 	ch := make(chan string)
// 	fn := func(i int) {
// 		ch <- vs[i](name)
// 	}
// 	for i, _ := range vs {
// 		go fn(i)
// 	}
// 	//runtime.Gosched()
// 	return <-ch
// }

// func main() {
// 	ret := exec("111", func(n string) string {
// 		return n + "func1"
// 	}, func(n string) string {
// 		return n + "func2"
// 	}, func(n string) string {
// 		return n + "func3"
// 	}, func(n string) string {
// 		return n + "func4"
// 	})
// 	fmt.Println(ret)
// }

func main() {
	var once sync.Once

	for i := 0; i <= 10; i++ {
		go once.Do(func() {
			fmt.Println("hello world")
		})
	}

	time.Sleep(time.Second * 2)
}
