package main

import (
	"fmt"
	"time"

	"github.com/angelbeltran/go-observer"
)

func runPublisher(prop observer.Property[int]) {
	val := prop.Value()
	for {
		time.Sleep(time.Second)
		// update property
		val++
		prop.Update(val)
	}
}

func runObserver(id int, prop observer.Property[int]) {
	stream := prop.Observe()

	for {
		val := stream.Value()
		fmt.Printf("Observer: %d, Value: %d\n", id, val)

		select {
		// wait for changes
		case <-stream.Changes():
			// advance to next value
			stream.Next()
		}
	}
}

func main() {
	// create a property with initial value
	prop := observer.NewProperty(1)

	// run 10 observers
	for i := 0; i < 10; i++ {
		go runObserver(i, prop)
	}

	// run one publisher
	go runPublisher(prop)

	// terminate program after 10 seconds
	time.Sleep(10 * time.Second)
}
