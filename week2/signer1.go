package main

import (
	"fmt"
	"strconv"
	"sync"
)

//  сюда писать код
func ExecutePipeline(jobs ...job) {
	var wg sync.WaitGroup
	in := make(chan interface{})
	wg.Add(len(jobs))
	for _, j := range jobs {
		out := make(chan interface{})
		go func(in, out chan interface{}, j job) {
			defer wg.Done()
			j(in, out)
			defer close(out)
		}(in, out, j)
		in = out
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	nCh := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex
	data := <-in
	newData := strconv.Itoa(data.(int))
	go func(mu *sync.Mutex) {
		mu.Lock()

		mu.Unlock()
	}(&mu)
	var res string
	res += DataSignerCrc32(newData) + "~" + DataSignerCrc32(DataSignerMd5(newData))
	fmt.Println(res)
	out <- res
}

func MultiHash(in, out chan interface{}) {
	data := <-in
	var res string
	for i := 0; i <= 5; i++ {
		res += DataSignerCrc32(strconv.Itoa(i) + data.(string))
		fmt.Println(res)
	}
	out <- res
}

func CombineResults(in, out chan interface{}) {
	s := []string{}
	for v := range in {
		s = append(s, v.(string), "_")
	}
	fmt.Println(s)
}
