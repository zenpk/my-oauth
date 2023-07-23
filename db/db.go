package db

import (
	"github.com/zenpk/safe-csv-db"
	"log"
)

var Users *scd.Table
var RefreshTokens *scd.Table
var Clients *scd.Table

func Init(prepared, done chan struct{}) error {
	var err error

	Users, err = scd.OpenTable("./users.csv")
	if err != nil {
		return err
	}
	defer Users.Close()
	go func() {
		if err := Users.ListenChange(); err != nil {
			log.Fatalln(err)
		}
	}()

	RefreshTokens, err = scd.OpenTable("./refresh_tokens.csv")
	if err != nil {
		return err
	}
	defer RefreshTokens.Close()
	go func() {
		if err := RefreshTokens.ListenChange(); err != nil {
			log.Fatalln(err)
		}
	}()

	Clients, err = scd.OpenTable("./clients.csv")
	if err != nil {
		return err
	}
	defer Clients.Close()
	go func() {
		if err := Clients.ListenChange(); err != nil {
			log.Fatalln(err)
		}
	}()

	prepared <- struct{}{}
	<-done
	return nil
}
