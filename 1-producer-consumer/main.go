//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

func producer(stream Stream, in chan<- *Tweet) {

	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(in)
			return
		}

		in <- tweet
	}
}

func consumer(in <-chan *Tweet, done chan<- bool) {

	for t := range in {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}

	done <- true

}

func main() {
	start := time.Now()

	in := make(chan *Tweet)

	done := make(chan bool)

	stream := GetMockStream()

	// Producer
	go producer(stream, in)

	// Consumer
	go consumer(in, done)

	<-done

	fmt.Printf("Process took %s\n", time.Since(start))
}
