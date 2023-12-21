package main

import (
	"fmt"
	"time"
)

func main() {
	cancelCh := make(chan struct{})
	valueCh := make(chan int)

	start := time.Now()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			v := time.Since(start)
			valueCh <- int(v)

			// emit cancel
			if time.Since(start) >= 5*time.Second {
				cancelCh <- struct{}{}
				return
			}
		}
	}()

	for {
		select {
		case <-cancelCh:
			fmt.Printf("closing channel")
			close(valueCh)
			return
		case v := <-valueCh:
			fmt.Printf("Second: %v\n", v/int(time.Second))
		}
	}

}
