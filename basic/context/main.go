package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Context Done ...")
					close(dst)
					return
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}
	// ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()
	for n := range gen(ctx) {
		fmt.Println(n)
		// if n == 5 {
		// 	break
		// }
	}

	ctx = context.WithValue(context.Background(), "color", "red")

	fmt.Println(ctx.Value("color").(string))

}
