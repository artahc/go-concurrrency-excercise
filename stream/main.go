package main

import (
	"context"
	"fmt"
)

type Stream[T any] struct {
	ch <-chan T
}

func (s *Stream[T]) Subscribe(f func(v T)) {
	for v := range s.ch {
		f(v)
	}
}

func NewStream[T any](ctx context.Context, items []T) *Stream[T] {
	ch := make(chan T)

	go func() {
		defer close(ch)
		for _, v := range items {
			select {
			case <-ctx.Done():
				return
			default:
				value := v
				ch <- value
			}
		}
	}()

	return &Stream[T]{
		ch: ch,
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	stream := NewStream[int](ctx, []int{1, 2, 3, 4, 5})

	stream.Subscribe(func(v int) {
		fmt.Printf("Subscribed 1: %v\n", v)
	})

	streamStr := NewStream[string](ctx, []string{"1", "2"})

	streamStr.Subscribe(func(v string) {
		if v == "1" {
			cancel()
		}
		fmt.Printf("Subscriber 2: %v\n", v)
	})
}
