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

	logFile, err := os.OpenFile(utils.Conf.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	done := make(chan struct{})
	prepared := make(chan struct{})
	exited := make(chan struct{})
	dbInstance := new(db.Db)
	go func() {
		if err := dbInstance.Init(prepared, done); err != nil {
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

	handlerInstance := handlers.Handler{Db: dbInstance}
	if err := handlers.StartListening(handlerInstance); err != nil {
		panic(err)
	}
}
