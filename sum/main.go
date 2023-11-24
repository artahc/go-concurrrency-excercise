package main

import (
	"fmt"
	"sync"
)

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	wg := sync.WaitGroup{}
	ch := make(chan int)

	sum := func(num []int) {
		defer wg.Done()
		res := 0
		for _, i := range num {
			res += i
		}
		ch <- res
	}

	wg.Add(2)
	go sum(nums[0:4])
	go sum(nums[5:9])

	go func() {
		wg.Wait()
		close(ch)
	}()

	res := 0
	for value := range ch {
		res += value
	}

	fmt.Printf("Result = %v\n", res)
}
