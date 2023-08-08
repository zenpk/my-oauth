package db

import (
	"github.com/zenpk/safe-csv-db"
	"log"
)

var UserCsv *scd.Csv
var RefreshTokenCsv *scd.Csv
var ClientCsv *scd.Csv

func Init(prepared, done chan struct{}) error {
	var err error

	UserCsv, err = scd.OpenCsv("./db/user.csv", User{})
	if err != nil {
		return err
	}
	defer UserCsv.Close()
	go func() {
		if err := UserCsv.ListenChange(); err != nil {
			log.Fatalln(err)
		}
	}()

	ClientCsv, err = scd.OpenCsv("./db/client.csv", Client{})
	if err != nil {
		return err
	}
	defer ClientCsv.Close()
	go func() {
		if err := ClientCsv.ListenChange(); err != nil {
			log.Fatalln(err)
		}
	}()

	RefreshTokenCsv, err = scd.OpenCsv("./db/refresh_token.csv", RefreshToken{})
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
