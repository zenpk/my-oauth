package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zenpk/my-oauth/dal"
	"github.com/zenpk/my-oauth/handler"
	"github.com/zenpk/my-oauth/service"
	"github.com/zenpk/my-oauth/token"
	"github.com/zenpk/my-oauth/util"
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
		log.Println("gracefully exited")
	}()

	conf := new(util.Configuration)
	if err := conf.Init(*mode); err != nil {
		panic(err)
	}

	logger := new(util.Logger)
	if err := logger.Init(conf); err != nil {
		panic(err)
	}
	defer func() {
		log.Println("exiting")
		if err := logger.Close(); err != nil {
			cleanUpErr = errors.Join(cleanUpErr, err)
		}
	}()

	db := new(dal.Database)
	if err := db.Init(conf); err != nil {
		panic(err)
	}

	authInfo := new(util.AuthorizationInfo)
	authInfo.Init(conf)
	tk := new(token.Token)
	if err := tk.Init(conf, logger); err != nil {
		panic(err)
	}
	service := new(service.Service)
	service.Init(conf, db)

	hd := new(handler.Handler)
	hd.Init(conf, logger, service, authInfo, tk)

	// clean up
	osSignalChan := make(chan os.Signal, 2)
	signal.Notify(osSignalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-osSignalChan
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := hd.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	log.Println("started")
	if err := hd.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
