// package main

// import (
// 	"fmt"
// 	"os"
// 	"strings"
// 	"sync"
// )

// func main() {
// 	data, err := os.ReadFile("sample.txt")
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 		return
// 	}

// 	mu := sync.Mutex{}
// 	wg := sync.WaitGroup{}
// 	maps := make(map[string]int)

// 	lines := strings.Split(string(data), "\n")
// 	countOccurences := func(line string) {
// 		defer wg.Done()
// 		words := strings.Fields(line)
// 		for _, word := range words {
// 			word = strings.ToLower(strings.Trim(word, ".,!?;:'\"()[]{}-"))
// 			if word != "" {
// 				mu.Lock()
// 				maps[word] += 1
// 				mu.Unlock()
// 			}
// 		}
// 	}

// 	for _, line := range lines {
// 		wg.Add(1)
// 		go countOccurences(line)
// 	}

// 	wg.Wait()

//		mu.Lock()
//		for key, count := range maps {
//			fmt.Printf("Word: %s -> %v\n", key, count)
//		}
//		mu.Unlock()
//	}
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	data, err := os.ReadFile("sample.txt")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	wordCountChan := make(chan map[string]int)

	for _, line := range lines {
		go countOccurrences(line, wordCountChan)
	}

	// Collect word counts
	wordCounts := make(map[string]int)
	for i := 0; i < len(lines); i++ {
		counts := <-wordCountChan
		for word, count := range counts {
			wordCounts[word] += count
		}
	}
	// Print the final result
	for key, count := range wordCounts {
		fmt.Printf("Word: %s -> %v\n", key, count)
	}
}

func countOccurrences(line string, wordCountChan chan map[string]int) {
	words := strings.Fields(line)
	counts := make(map[string]int)
	for _, word := range words {
		word = strings.ToLower(strings.Trim(word, ".,!?;:'\"()[]{}-"))
		if word != "" {
			counts[word]++
		}
	}
	wordCountChan <- counts
}
