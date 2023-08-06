package main

import (
	"flag"
	"fmt"
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
	var str string
	str, _ = utils.RandString(10)
	fmt.Println(str)
	str, _ = utils.RandString(10)
	fmt.Println(str)
	str, _ = utils.RandString(10)
	fmt.Println(str)
	str, _ = utils.RandString(10)
	fmt.Println(str)

	if err := handlers.StartListening(); err != nil {
		log.Fatalln(err)
	}
}
