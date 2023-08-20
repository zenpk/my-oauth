package main

import (
	"flag"
	"github.com/zenpk/my-oauth/db"
	"github.com/zenpk/my-oauth/handlers"
	"github.com/zenpk/my-oauth/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	mode = flag.String("mode", "dev", "define program mode")
)

func main() {
	flag.Parse()
	if err := utils.Init(*mode); err != nil {
		panic(err)
	}

	done := make(chan struct{})
	prepared := make(chan struct{})
	exited := make(chan struct{})
	go func() {
		if err := db.Init(prepared, done); err != nil {
			panic(err)
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
		log.Println("gracefully exited")
		os.Exit(0)
	}()

	if err := handlers.StartListening(); err != nil {
		panic(err)
	}
}
