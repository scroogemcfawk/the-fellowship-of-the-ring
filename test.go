package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	var n = 4

	chans := make([]chan string, n)
	for i := range chans {
		chans[i] = make(chan string, 2)
	}

	for i := 0; i < n-1; i++ {
		wg.Add(1)
		go func(num int) {
			sendMsg(chans[getNext(num, n)], num)
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(num int) {
			readMsg(chans[num], num)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func sendMsg(ch chan string, cur int) {
	ch <- fmt.Sprintf("A message from %d", cur)
}

func readMsg(ch chan string, cur int) {
	select {
	case msg := <-ch:
		fmt.Printf("%d got message: %s\n", cur, msg)
	default:
		return
		//fmt.Printf("%d got no message\n", cur)
	}
}

//func getNext(c, n int) int {
//	if c == n-1 {
//		return 0
//	} else {
//		return c + 1
//	}
//}

//func getPrev(c, n int) int {
//	if c == 0 {
//		return n - 1
//	} else {
//		return c - 1
//	}
//}
