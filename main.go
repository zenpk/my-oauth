package main

import (
	"flag"
	"github.com/zenpk/my-oauth/db"
	"github.com/zenpk/my-oauth/handlers"
	"github.com/zenpk/my-oauth/utils"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	mode = flag.String("mode", "dev", "define program mode")
)

func main() {
	// catch the panic before exit
	defer func() {
		if err := recover(); err != nil {
			printStack()
			panic(err)
		}
	}()

	flag.Parse()
	if err := utils.Init(*mode); err != nil {
		log.Fatalln(err)
	}

	done := make(chan struct{})
	prepared := make(chan struct{})
	exited := make(chan struct{})
	go func() {
		if err := db.Init(prepared, done); err != nil {
			log.Fatalln(err)
		}
		exited <- struct{}{}
	}()
	<-prepared

	// clean up
	osSignalChan := make(chan os.Signal, 2)
	signal.Notify(osSignalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-osSignalChan
		done <- struct{}{}
		<-exited
		log.Fatalln("gracefully exited")
	}()

	if err := handlers.StartListening(); err != nil {
		log.Fatalln(err)
	}
}

func printStack() {
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, true)
	log.Printf("=== BEGIN Goroutine stack trace ===\n%s\n=== END Goroutine stack trace ===", buf[:stackSize])
}
