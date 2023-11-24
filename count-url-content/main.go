package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

func fetch(urls ...string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, url := range urls {
			res, err := http.Get(url)
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}
			fmt.Printf("%s\n", data)
			out <- string(data)
		}
	}()

	return out
}

func countWords(ch <-chan string) <-chan map[string]int {
	out := make(chan map[string]int)

	maps := make(map[string]int)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for c := range ch {
			sentences := strings.Split(c, ".")
			wg.Add(len(sentences))

			countWordsInSentence := func(sentence string) {
				defer wg.Done()
				words := strings.Fields(sentence)
				for _, word := range words {
					word = strings.ToLower(strings.Trim(word, ".,!?;:'\"()[]{}-"))
					mu.Lock()
					maps[word]++
					mu.Unlock()
				}
			}

			for _, sentence := range sentences {
				go countWordsInSentence(sentence)
			}
		}
	}()

	go func() {
		wg.Wait()
		out <- maps
		close(out)
	}()

	return out
}

func main() {
	urls := []string{
		"https://baconipsum.com/api/?type=meat-and-filler&paras=5&format=text",
		// "https://baconipsum.com/api/?type=all-meat&paras=3&format=text",
		// "https://baconipsum.com/api/?type=all-meat&sentences=150&format=text",
	}

	ch := fetch(urls...)
	countWords := countWords(ch)

	counts := make(map[string]int)
	for key, value := range <-countWords {
		counts[key] += value
	}

	for key, value := range counts {
		fmt.Printf("%s: %v\n", key, value)
	}
}
