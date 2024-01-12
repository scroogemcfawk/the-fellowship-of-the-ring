package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Token struct {
	data      string
	recipient int
	ttl       int
}

var wg sync.WaitGroup

func main() {

	// read user input - number of workers
	var n int
	fmt.Scanf("%d", &n)

	// create channels
	// the ith channels message is a message to ith worker from ith - 1
	chans := make([]chan Token, n)
	for i := range chans {
		// making buffered channels so they don't block the coroutine
		chans[i] = make(chan Token, 2)
	}

	// create message
	var message Token
	message.data = "One ring to rule them all, one ring to find them, one ring to bring them all and in the darkness bind them."
	message.ttl = rand.Intn(n)
	message.recipient = rand.Intn(n)
	//message.recipient = 1

	println()
	printToken(message)
	println()

	// send message to 1st channel
	chans[0] <- message

	for x := 0; x < n; x++ {
		wg.Add(1)
		go func(x int) {
			process(chans, x, n)
			wg.Done()
		}(x)
		time.Sleep(time.Microsecond)
	}

	wg.Wait()
}

func process(chans []chan Token, current int, n int) {
	// ignore closing closed channel
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	for {
		select {
		// read msg
		case msg, open := <-chans[current]:

			next := getNext(current, n)
			if !open {
				close(chans[next])
				return
			}

			// print or send msg
			if msg.recipient == current {
				fmt.Printf("%d: %s\n", current, msg.data)
				close(chans[next])
			} else {
				if msg.ttl > 0 {
					msg.ttl -= 1
					chans[next] <- msg
					fmt.Printf("%d: Sending to %d\n", current, next)
				} else {
					fmt.Printf("%d: Message for recipient %d died\n", current, msg.recipient)
					close(chans[next])
				}
			}

		default:
			//fmt.Printf("No msg\n")
			continue
		}
	}
}

func getNext(c, n int) int {
	if c == n-1 {
		return 0
	} else {
		return c + 1
	}
}

func printToken(t Token) {
	fmt.Printf("Token:\nMessage: %s\nRecipient: %d\nttl: %d\n", t.data, t.recipient, t.ttl)
}
