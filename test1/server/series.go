package server

import "fmt"

func init(){
	fmt.Println("init1")
}

func Fib(n int) []int {

	fib := []int{1, 1}
	for i := 2; i < n; i++ {
		fib = append(fib, fib[i-2]+fib[i-1])
	}
	return fib
}
