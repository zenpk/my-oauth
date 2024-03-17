package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zenpk/my-oauth/dal"
	"github.com/zenpk/my-oauth/handlers"
	"github.com/zenpk/my-oauth/utils"
)

var mode = flag.String("mode", "dev", "define program mode")

func main() {
	flag.Parse()
	// graceful exit
	var cleanUpErr error
	defer func() {
		if cleanUpErr != nil {
			panic(cleanUpErr)
		}
		fmt.Println("gracefully exited") // need to use fmt because at this point the logFile is already closed
	}()

	conf := new(utils.Configuration)
	if err := conf.Init(*mode); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(conf.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(logFile)
	defer func() {
		log.Println("exiting")
		if err := logFile.Close(); err != nil {
			cleanUpErr = errors.Join(cleanUpErr, err)
		}
	}()

	stop := make(chan struct{})
	preparing := make(chan struct{})
	dbRunning := make(chan struct{})
	db := new(dal.Db)
	go func() {
		if err := db.Init(preparing, stop); err != nil {
			panic(err)
		}
		close(dbRunning)
	}()
	<-preparing

	// clean up
	osSignalChan := make(chan os.Signal, 2)
	signal.Notify(osSignalChan, os.Interrupt, syscall.SIGTERM)

	handlerInstance := handlers.Handler{Db: db}
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
