package main

import (
	"fmt"
	"sort"
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
	var wg sync.WaitGroup
	var mu sync.Mutex
	for v := range in {
		data := strconv.Itoa(v.(int))
		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			ch1 := make(chan string)
			ch2 := make(chan string)
			go func() {
				ch1 <- DataSignerCrc32(data)
			}()
			go func() {
				mu.Lock()
				mData := DataSignerMd5(data)
				mu.Unlock()
				ch2 <- DataSignerCrc32(mData)
			}()
			res := <-ch1 + "~" + <-ch2
			fmt.Println(res)
			out <- res
		}(data)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	// var mu sync.Mutex

	var wg sync.WaitGroup
	for v := range in {
		data := v.(string)
		// fmt.Println(data)
		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			var wgA sync.WaitGroup
			res := make([]string, 6)
			for i := 0; i <= 5; i++ {
				wgA.Add(1)
				go func(i int) {
					defer wgA.Done()
					res[i] = DataSignerCrc32(strconv.Itoa(i) + data)
				}(i)
			}
			wgA.Wait()
			var ress string
			for _, v := range res {
				ress += v
			}
			fmt.Println(ress)
			out <- ress
		}(data)
	}
	wg.Wait()
}

// func MultiHash(in, out chan interface{}) {
// 	var mu sync.Mutex
// 	var wg sync.WaitGroup
// 	for v := range in {
// 		data := v.(string)
// 		fmt.Println(data)
// 		wg.Add(1)
// 		var res string
// 		go func(wg *sync.WaitGroup, data string, res string, mu *sync.Mutex) {
// 			defer wg.Done()
// 			for i := 0; i <= 5; i++ {
// 				wg.Add(1)
// 				go func(i int) {
// 					defer wg.Done()
// 					mu.Lock()
// 					res += DataSignerCrc32(strconv.Itoa(i) + data)
// 					fmt.Println(res)
// 					mu.Unlock()
// 				}(i)
// 			}
// 		}(&wg, data, res, &mu)
// 		wg.Wait()
// 		fmt.Println(res)
// 		out <- res
// 	}
// }

func CombineResults(in, out chan interface{}) {
	s := []string{}
	for v := range in {
		s = append(s, v.(string))
	}
	// fmt.Println(s)
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
	var res string
	for i, v := range s {
		if i == len(s)-1 {
			res += v
		} else {
			res += v + "_"
		}
	}
	out <- res
}
