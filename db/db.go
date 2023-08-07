package db

import (
	"github.com/zenpk/safe-csv-db"
	"log"
)

var UserCsv *scd.Table
var RefreshTokenCsv *scd.Table
var ClientCsv *scd.Table

func Init(prepared, done chan struct{}) error {
	var err error

	UserCsv, err = scd.OpenTable("./db/user.csv")
	if err != nil {
		return err
	}
	defer UserCsv.Close()
	go func() {
		if err := UserCsv.ListenChange(); err != nil {
			log.Fatalln(err)
		}
	}()

	ClientCsv, err = scd.OpenTable("./db/client.csv")
	if err != nil {
		return err
	}
	defer ClientCsv.Close()
	go func() {
		if err := ClientCsv.ListenChange(); err != nil {
			log.Fatalln(err)
		}
	}()

	RefreshTokenCsv, err = scd.OpenTable("./db/refresh_token.csv")
	if err != nil {
		return err
	}
	defer RefreshTokenCsv.Close()
	go func() {
		if err := RefreshTokenCsv.ListenChange(); err != nil {
			log.Fatalln(err)
		}
	}()

	prepared <- struct{}{}
	<-done
	return nil
}
