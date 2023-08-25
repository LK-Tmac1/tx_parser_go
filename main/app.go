package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	initApp()
}

func initApp() {
	fmt.Println("initializing app...")

	p := NewParser(NewInMemoryStore())
	if err := p.initBlock(); err != nil {
		log.Fatal("failed to initialize block info...")
	}

	// run a background job to update transactions on the fly
	fmt.Println("initializing background update...")
	go p.initBackgroundUpdate(5 * time.Second)

	fmt.Println("starting service on 8080...")
	initHandler("8080", p)
}
