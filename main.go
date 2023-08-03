package main

import (
	"github.com/zenpk/my-oauth/db"
	"github.com/zenpk/my-oauth/handlers"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
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
