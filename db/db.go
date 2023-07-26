package db

import (
	"github.com/zenpk/safe-csv-db"
	"log"
)

var UserTable *scd.Table
var RefreshTokenTable *scd.Table
var ClientTable *scd.Table

func Init(prepared, done chan struct{}) error {
	var err error

	UserTable, err = scd.OpenTable("./user.csv")
	if err != nil {
		return err
	}
	defer UserTable.Close()
	go func() {
		if err := UserTable.ListenChange(); err != nil {
			log.Fatalln(err)
		}
	}()

	ClientTable, err = scd.OpenTable("./client.csv")
	if err != nil {
		return err
	}
	defer ClientTable.Close()
	go func() {
		if err := ClientTable.ListenChange(); err != nil {
			log.Fatalln(err)
		}
	}()

	RefreshTokenTable, err = scd.OpenTable("./refresh_token.csv")
	if err != nil {
		return err
	}
	defer RefreshTokenTable.Close()
	go func() {
		if err := RefreshTokenTable.ListenChange(); err != nil {
			log.Fatalln(err)
		}
	}()

	prepared <- struct{}{}
	<-done
	return nil
}
