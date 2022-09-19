package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// var urls = flag.String("urls", "", "Список URL через запятую")

func main() {
	var mRoutines = flag.Int("max-routines", 5, "Максимальное количество горутин")
	flag.Parse()
	queue := make(chan bool, *mRoutines)
	var counters sync.Map
	var wg sync.WaitGroup
	urls := []string {"https://google.com", "https://habr.com", "https://stackoverflow.com", "https://vk.com", "https://youtube.com", "https://go.dev", "https://google.com"}
	
	wg.Add(len(urls))
	for _, url := range urls {
		queue <- true
		url := url
		go worker(queue, &wg, func ()  {
			parseUrl(url, &counters)
		})
	}

	wg.Wait()
	close(queue)

	var total int = 0;

	counters.Range(func(key, value any) bool {
		fmt.Printf("Count for %s: %d \n", key, value)
		total += value.(int)
		return true
	})

	fmt.Printf("Total: %d\n", total)
}

func worker(queue <-chan bool, wg *sync.WaitGroup, handler func()) {
	handler()
	wg.Done()
	<-queue
}

func parseUrl(url string, counters *sync.Map) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Ошибка обработки %s %s \n", url, err)

		return
	}

	scanner := bufio.NewScanner(resp.Body)

	var count int = 0;

	for scanner.Scan() {
		count += strings.Count(scanner.Text(), "Go")
	}

	counters.Store(url, count)
}