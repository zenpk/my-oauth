package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/zenpk/my-oauth/db"
	"github.com/zenpk/my-oauth/handlers"
	"github.com/zenpk/my-oauth/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	mode = flag.String("mode", "dev", "define program mode")
)

func main() {
	flag.Parse()
	// graceful exit
	var cleanUpErrors []error
	defer func() {
		for _, err := range cleanUpErrors {
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("gracefully exited") // need to use fmt because at this point the logFile is already closed
	}()

	if err := utils.Init(*mode); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(utils.Conf.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(logFile)
	defer func() {
		log.Println("exiting")
		if err := logFile.Close(); err != nil {
			cleanUpErrors = append(cleanUpErrors, err)
		}
	}()

	stop := make(chan struct{})
	preparing := make(chan struct{})
	dbRunning := make(chan struct{})
	dbInstance := new(db.Db)
	go func() {
		if err := dbInstance.Init(preparing, stop); err != nil {
			panic(err)
		}
		close(dbRunning)
	}()
	<-preparing

	// clean up
	osSignalChan := make(chan os.Signal, 2)
	signal.Notify(osSignalChan, os.Interrupt, syscall.SIGTERM)

	handlerInstance := handlers.Handler{Db: dbInstance}
	server := handlers.CreateServer(handlerInstance)

	go func() {
		<-osSignalChan
		close(stop)
		<-dbRunning
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	log.Printf("start listening at %v\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
